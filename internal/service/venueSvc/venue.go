package venueSvc

import (
	"context"
	"fmt"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
)

// CreateVenue 创建新场地
func CreateVenue(ctx context.Context, venue *models.Venue) (int, error) {
	// 校验必填字段
	if venue.Name == "" {
		return 0, ErrVenueNameRequired
	}
	if venue.BuildingID == 0 {
		return 0, ErrVenueBuildingRequired
	}
	normalizedEquipments, err := normalizeVenueEquipments(venue.EquipmentsRaw)
	if err != nil {
		return 0, err
	}
	venue.EquipmentsRaw = normalizedEquipments

	// 创建时默认设置为启用状态
	venue.IsActive = true

	// 创建场地
	venueID, err := repository.CreateNewVenue(venue)
	if err != nil { // ENH: 此处未区分数据库错误与外链错误，可以根据需要进一步细化错误类型
		return 0, fmt.Errorf("%w: %w", ErrVenueCreateInDB, err)
	}

	return venueID, nil
}

// UpdateVenue 更新场地信息
func UpdateVenue(ctx context.Context, updates *models.Venue) error {
	existingVenue, err := repository.GetVenueByID(updates.ID)
	if err != nil {
		return err
	}

	normalizedEquipments, err := normalizeVenueEquipments(updates.EquipmentsRaw)
	if err != nil {
		return err
	}
	updates.EquipmentsRaw = normalizedEquipments

	// 应用更新（只更新非零值字段）
	utils.UpdateField(&existingVenue.Name, updates.Name)
	utils.UpdateField(&existingVenue.BuildingID, updates.BuildingID)
	utils.UpdateField(&existingVenue.TypeID, updates.TypeID)
	utils.UpdateField(&existingVenue.Capacity, updates.Capacity)
	utils.UpdateField(&existingVenue.Description, updates.Description)
	utils.UpdateField(&existingVenue.CoverImageToken, updates.CoverImageToken)
	if len(updates.EquipmentsRaw) > 0 {
		existingVenue.EquipmentsRaw = updates.EquipmentsRaw
	}

	// 如果更新了楼区ID，校验其是否存在
	if updates.BuildingID != 0 {
		if err := checkBuildingExists(updates.BuildingID); err != nil {
			return err
		}
	}

	// 执行更新
	if err := repository.UpdateVenue(existingVenue); err != nil {
		return fmt.Errorf("%w: %w", ErrVenueUpdateInDB, err)
	}

	return nil
}

// DeleteVenue 删除场地（软删除）
func DeleteVenue(ctx context.Context, venueID int) error {
	// 执行软删除
	if err := repository.DeleteVenueByID(venueID); err != nil {
		return fmt.Errorf("%w: %w", ErrVenueDeleteInDB, err)
	}

	return nil
}

// VenueListOptions 场地列表查询选项
// 非数据模型，仅用于参数传递
type VenueListOptions struct {
	BuildingIDs []int
	TypeIDs     []int
	Permissions []venuePermission.VenuePerm
	SearchQuery string
	Offset      int
	Limit       int
	VAGID       int
	SysPerm     uint64
}

// ListVenues 列出场地
func ListVenues(ctx context.Context, opts VenueListOptions) ([]*models.Venue, error) {
	// 设置默认分页大小
	if opts.Limit == 0 {
		opts.Limit = 12
	}

	// 查询场地列表
	venues, err := repository.ListVenuesWithQuery(repository.VenueQueryOptions{
		BuildingIDs: opts.BuildingIDs,
		TypeIDs:     opts.TypeIDs,
		Permissions: opts.Permissions,
		SearchQuery: opts.SearchQuery,
		Offset:      opts.Offset,
		Limit:       opts.Limit,
		VAGID:       opts.VAGID,
		SysPerm:     opts.SysPerm,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrVenueQueryInDB, err)
	}

	return venues, nil
}

// GetAccessibleBuildingsAndCampuses 根据权限组获取可访问的楼区和校区信息
func GetAccessibleBuildingsAndCampuses(ctx context.Context, vagid int, allowAll bool) ([]*models.VenueBuilding, []*models.VenueCampus, error) {
	buildings, err := repository.GetAccessibleBuildingList(vagid, allowAll)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrBuildingQueryInDB, err)
	}
	campuses, err := repository.GetAccessibleCampusList(vagid, allowAll)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrCampusQueryInDB, err)
	}
	return buildings, campuses, nil
}

// func checkVenueExists(venueID int) error {
// 	exists, err := repository.VenueExists(venueID)
// 	if err != nil {
// 		return fmt.Errorf("%w: %w", ErrVenueNotFound, err)
// 	}
// 	if !exists {
// 		return ErrVenueNotFound
// 	}
// 	return nil
// }

// checkBuildingExists 检查楼区是否存在
func checkBuildingExists(buildingID int) error {
	exists, err := repository.VenueBuildingExists(buildingID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrVenueBuildingInvalid, err)
	}
	if !exists {
		return ErrVenueBuildingInvalid
	}
	return nil
}
