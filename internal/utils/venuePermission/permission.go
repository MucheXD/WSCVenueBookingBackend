package venuePermission

import (
	"log/slog"
	"sync"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

// venueAccessCache 存储权限组的内存缓存，key为VAGID，value为VenueAccess指针
var (
	venueAccessCache = make(map[int]*models.VenueAccess)
	cacheLock        sync.RWMutex
)

// DataLoader 定义数据加载器接口，用于获取权限数据
type DataLoader interface {
	GetAllVenueAccessGroupIDs() ([]int, error)
	GetVenueAccessGroupByID(vagid int) (*models.VenueAccess, error)
}

// RefreshVenueAccessCache 从数据库加载所有权限数据到内存缓存
// 这个函数应该在应用启动时和更新权限后调用
// 由于 util 不应该依赖 repo，因此使用依赖注入的方式获取数据加载器
func RefreshVenueAccessCache(loader DataLoader) error {
	slog.Debug("Starting to refresh venue access cache from database")

	// 获取所有权限组ID
	vagids, err := loader.GetAllVenueAccessGroupIDs()
	if err != nil {
		slog.Error("Failed to get all venue access group IDs", "error", err)
		return err
	}

	// 临时缓存，避免加载失败时污染现有缓存
	tempCache := make(map[int]*models.VenueAccess)

	// 加载每个权限组
	for _, vagid := range vagids {
		va, err := loader.GetVenueAccessGroupByID(vagid)
		if err != nil {
			slog.Warn("Failed to load venue access group", "vagid", vagid, "error", err)
			continue
		}
		if va != nil {
			tempCache[vagid] = va
		}
	}

	// 原子性地替换缓存
	cacheLock.Lock()
	defer cacheLock.Unlock()
	venueAccessCache = tempCache

	slog.Info("Venue access cache refreshed successfully", "count", len(venueAccessCache))
	return nil
}

// GetVenueAccessByVAGID 从缓存中获取指定权限组的权限信息
func GetVenueAccessByVAGID(vagid int) *models.VenueAccess {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	return venueAccessCache[vagid]
}

// CheckVenuePermission 检查用户对特定场地的权限
// 返回值：
// - true: 用户拥有指定权限
// - false: 用户不拥有指定权限
func CheckVenuePermission(vagid int, venueID int, perm VenuePerm) bool {
	va := GetVenueAccessByVAGID(vagid)
	if va == nil {
		return false
	}

	switch perm {
	case Reserve:
		return va.HasReserve(venueID)
	case Approval:
		return va.HasApproval(venueID)
	case Edit:
		return va.HasEdit(venueID)
	case Manage:
		return va.HasManage(venueID)
	default:
		return false
	}
}

// InvalidateCache 清空权限缓存（通常在权限更新后调用）
func InvalidateCache() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	venueAccessCache = make(map[int]*models.VenueAccess)
	slog.Debug("Venue access cache invalidated")
}
