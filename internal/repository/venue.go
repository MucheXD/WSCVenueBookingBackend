package repository

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"gorm.io/gorm"
)

type VenueEntity struct {
	VenueID         int            `gorm:"column:venue_id;primaryKey"`
	Name            string         `gorm:"column:name_text"`
	BuildingID      int            `gorm:"column:location_building_id"`
	TypeID          int            `gorm:"column:type_id"`
	Capacity        int            `gorm:"column:capacity"`
	Description     string         `gorm:"column:description_text"`
	CoverImageToken *string        `gorm:"column:cover_image_token"`
	EquipmentsRaw   []byte         `gorm:"column:equipments"`
	IsActive        bool           `gorm:"column:is_active"`
	DeletedAt       gorm.DeletedAt `gorm:"column:delete_at"`
}

func (VenueEntity) TableName() string {
	return "venues"
}

func CreateNewVenue(modelV *models.Venue) (int, error) {
	var venueEntity VenueEntity
	venueEntity.fromDomain(modelV)
	if err := database.DB.Create(&venueEntity).Error; err != nil {
		return 0, err
	}
	return venueEntity.VenueID, nil
}

func GetVenueByID(venueID int) (*models.Venue, error) {
	var venueEntity VenueEntity
	txDB := database.DB.
		Model(&VenueEntity{}).
		Where(&VenueEntity{VenueID: venueID}).
		Take(&venueEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return venueEntity.toDomain(), nil
}

func ListVenueBodies() ([]*models.Venue, error) {
	var venueEntities []VenueEntity
	if err := database.DB.
		Model(&VenueEntity{}).
		Order("venue_id ASC").
		Find(&venueEntities).Error; err != nil {
		return nil, err
	}

	venues := make([]*models.Venue, 0, len(venueEntities))
	for _, entity := range venueEntities {
		venues = append(venues, entity.toDomain())
	}
	return venues, nil
}

// VenueExists 检查场地是否存在
func VenueExists(venueID int) (bool, error) {
	var count int64
	txDB := database.DB.
		Model(&VenueEntity{}).
		Where(&VenueEntity{VenueID: venueID}).
		Count(&count)
	if txDB.Error != nil {
		return false, txDB.Error
	}
	return count > 0, nil
}

func UpdateVenue(modelV *models.Venue) error {
	var venueEntity VenueEntity
	venueEntity.fromDomain(modelV)
	if err := database.DB.Model(&VenueEntity{}).
		Where(&VenueEntity{VenueID: venueEntity.VenueID}).
		Updates(&venueEntity).Error; err != nil {
		return err
	}
	return nil
}

func DeleteVenueByID(venueID int) error {
	// 使用软删除
	txDB := database.DB.
		Model(&VenueEntity{}).
		Where(&VenueEntity{VenueID: venueID}).
		Delete(&VenueEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// VenueQueryOptions 场地查询选项
type VenueQueryOptions struct {
	BuildingIDs []int                       // 楼区ID筛选
	TypeIDs     []int                       // 类型ID筛选
	Permissions []venuePermission.VenuePerm // 权限筛选
	SearchQuery string                      // 搜索关键词
	Offset      int                         // 分页偏移
	Limit       int                         // 单页大小
	VAGID       int                         // 用户的权限组ID
	SysPerm     uint64                      // 用户的系统权限
}

// ListVenuesWithQuery 列出满足权限条件的场地
func ListVenuesWithQuery(opts VenueQueryOptions) ([]*models.Venue, error) {

	// SQL设计：使用JOIN联表查询venue_accesses，避免N+1查询问题
	// 目标 SQL 指令如下：
	// SELECT DISTINCT v.* FROM venues v
	// INNER JOIN venue_accesses va ON v.venue_id = va.venue_id
	// WHERE va.vagid = ? AND (va.allow_reserve = 1 OR va.allow_approval = 1 OR ...)
	// AND v.location_building_id IN (...)
	// AND v.type_id IN (...)
	// AND v.name_text LIKE ?
	// LIMIT ? OFFSET ?

	query := database.DB.Model(&VenueEntity{}).
		Distinct("venues.*")

	hasAllReserve := systemPermission.SatisfyAny(opts.SysPerm, systemPermission.AllVenueReservation)
	hasAllApproval := systemPermission.SatisfyAny(opts.SysPerm, systemPermission.AllVenueApproval)
	hasAllManage := systemPermission.SatisfyAny(opts.SysPerm, systemPermission.AllVenueManage)
	hasAllEdit := systemPermission.SatisfyAny(opts.SysPerm, systemPermission.AllVenueEdit)
	hasAnyAllVenuePerm := hasAllReserve || hasAllApproval || hasAllManage || hasAllEdit

	// 权限筛选：
	// 1) 无权限筛选器时：若无任意AllVenue...系统权限，则按venue_accesses + vagid限制；有则可列出全部。
	// 2) 有权限筛选器时：对应AllVenue...系统权限可直接满足；其余权限需通过venue_accesses + vagid限制。
	if len(opts.Permissions) == 0 {
		if !hasAnyAllVenuePerm {
			query = query.Joins("INNER JOIN venue_accesses va ON venues.venue_id = va.venue_id")
			query = query.Where("va.vagid = ?", opts.VAGID)
		}
	} else {
		permConditions := database.DB.Where("1 = 0")
		needVenueAccessFilter := false

		for _, perm := range opts.Permissions {
			switch perm {
			case venuePermission.Reserve:
				if !hasAllReserve {
					permConditions = permConditions.Or("va.allow_reserve = ?", true)
					needVenueAccessFilter = true
				}
			case venuePermission.Approval:
				if !hasAllApproval {
					permConditions = permConditions.Or("va.allow_approval = ?", true)
					needVenueAccessFilter = true
				}
			case venuePermission.Manage:
				if !hasAllManage {
					permConditions = permConditions.Or("va.allow_manage = ?", true)
					needVenueAccessFilter = true
				}
			case venuePermission.Edit:
				if !hasAllEdit {
					permConditions = permConditions.Or("va.allow_edit = ?", true)
					needVenueAccessFilter = true
				}
			}
		}

		if needVenueAccessFilter {
			query = query.Joins("INNER JOIN venue_accesses va ON venues.venue_id = va.venue_id")
			query = query.Where("va.vagid = ?", opts.VAGID)
			query = query.Where(permConditions)
		}
	}

	// 楼区筛选
	if len(opts.BuildingIDs) > 0 {
		query = query.Where("venues.location_building_id IN ?", opts.BuildingIDs)
	}

	// 类型筛选
	if len(opts.TypeIDs) > 0 {
		query = query.Where("venues.type_id IN ?", opts.TypeIDs)
	}

	// 搜索关键词
	if opts.SearchQuery != "" {
		searchPattern := "%" + opts.SearchQuery + "%"
		query = query.Where("venues.name_text LIKE ? OR venues.description_text LIKE ?",
			searchPattern, searchPattern)
	}

	// 分页
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}

	var venueEntities []VenueEntity
	if err := query.Find(&venueEntities).Error; err != nil {
		return nil, err
	}

	venues := make([]*models.Venue, 0, len(venueEntities))
	for _, entity := range venueEntities {
		venues = append(venues, entity.toDomain())
	}
	return venues, nil
}

// GetVenueIDsByVAGID 根据VAGID获取所有有权限的场地ID列表
// 用于获取用户可访问的场地位置信息
// SQL: SELECT DISTINCT venue_id FROM venue_accesses WHERE vagid = ?
func GetVenueIDsByVAGID(vagid int) ([]int, error) {
	var venueIDs []int
	if err := database.DB.Model(&VenueAccessEntity{}).
		Where("vagid = ?", vagid).
		Distinct("venue_id").
		Pluck("venue_id", &venueIDs).Error; err != nil {
		return nil, err
	}
	return venueIDs, nil
}

// 联表查询可访问的楼区信息
func GetAccessibleBuildingList(vagid int, allowAll bool) ([]*models.VenueBuilding, error) {
	var buildingEntities []VenueBuildingEntity

	// SQL设计：使用JOIN联表查询满足“条件A”的 venue_buildings 条目
	// INNER JOIN 实现了嵌套查询的效果，使用以下逻辑：
	// 条件A：该 venue_buildings 条目被满足“条件B”的 venues 条目引用
	// 条件B：该 venues 条目被满足“条件C”的 venue_accesses 条目引用
	// 条件C：该 venue_accesses 条目满足 vagid = ? 的条件
	// 每个条目全部满足条件任意次数后，只返回一次

	query := database.DB.Table("venue_buildings b")

	if !allowAll {
		query = query.
			Select("DISTINCT b.building_id, b.building_name, b.location_campus_id").
			Joins("INNER JOIN venues v ON b.building_id = v.location_building_id").
			Joins("INNER JOIN venue_accesses va ON v.venue_id = va.venue_id").
			Where("va.vagid = ?", vagid)
	}

	err := query.Scan(&buildingEntities).Error
	if err != nil {
		return nil, err
	}
	buildings := make([]*models.VenueBuilding, 0, len(buildingEntities))
	for _, entity := range buildingEntities {
		buildings = append(buildings, entity.toDomain())
	}
	return buildings, nil
}

// 联表查询可访问的校区信息
func GetAccessibleCampusList(vagid int, allowAll bool) ([]*models.VenueCampus, error) {
	var campusEntities []VenueCampusEntity

	// SQL设计：使用JOIN联表查询满足“条件A”的 venue_campuses 条目
	// INNER JOIN 实现了嵌套查询的效果，使用以下逻辑：
	// 条件A：该 venue_campuses 条目被满足“条件B”的 venue_buildings 条目引用
	// 条件B：该 venue_buildings 条目被满足“条件C”的 venues 条目引用
	// 条件C：该 venues 条目被满足“条件D”的 venue_accesses 条目引用
	// 条件D：该 venue_accesses 条目满足 vagid = ? 的条件
	// 每个条目全部满足条件任意次数后，只返回一次

	query := database.DB.Table("venue_campuses c")

	if !allowAll {
		query = query.
			Select("DISTINCT c.campus_id, c.campus_name").
			Joins("INNER JOIN venue_buildings b ON c.campus_id = b.location_campus_id").
			Joins("INNER JOIN venues v ON b.building_id = v.location_building_id").
			Joins("INNER JOIN venue_accesses va ON v.venue_id = va.venue_id").
			Where("va.vagid = ?", vagid)
	}

	err := query.Scan(&campusEntities).Error
	if err != nil {
		return nil, err
	}
	campuses := make([]*models.VenueCampus, 0, len(campusEntities))
	for _, entity := range campusEntities {
		campuses = append(campuses, entity.toDomain())
	}
	return campuses, nil
}

// GetVenueTimeslotsInRange 获取指定场地的时间段占用情况
func GetVenueTimeslotsInRange(venueID int, rangeBegin, rangeEnd time.Time, limit int) ([]models.VenueTimeslot, error) {

	// SQL: SELECT * FROM venue_timeslots WHERE venue_id = ? AND start_time >= ? AND start_time <= ? ORDER BY start_time ASC LIMIT ?

	var entities []VenueTimeslotEntity
	query := database.DB.Model(&VenueTimeslotEntity{}).
		Where("venue_id = ?", venueID).
		Where("start_time >= ?", rangeBegin).
		Where("start_time <= ?", rangeEnd).
		Order("start_time ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	timeslots := make([]models.VenueTimeslot, 0, len(entities))
	for _, entity := range entities {
		timeslots = append(timeslots, *entity.toDomain())
	}
	return timeslots, nil
}

// Entity-Domain Conversion (Private Methods)

func (v *VenueEntity) fromDomain(modelV *models.Venue) {
	v.VenueID = modelV.ID
	v.Name = modelV.Name
	v.BuildingID = modelV.BuildingID
	v.TypeID = modelV.TypeID
	v.Capacity = modelV.Capacity
	v.Description = modelV.Description
	if modelV.CoverImageToken == "" {
		v.CoverImageToken = nil
	} else {
		v.CoverImageToken = &modelV.CoverImageToken
	}
	v.EquipmentsRaw = modelV.EquipmentsRaw
	v.IsActive = modelV.IsActive
}

func (v *VenueEntity) toDomain() *models.Venue {
	return &models.Venue{
		ID:              v.VenueID,
		Name:            v.Name,
		BuildingID:      v.BuildingID,
		TypeID:          v.TypeID,
		Capacity:        v.Capacity,
		Description:     v.Description,
		CoverImageToken: defaultImageIfEmpty(v.CoverImageToken),
		EquipmentsRaw:   v.EquipmentsRaw,
		IsActive:        v.IsActive,
	}
}

func defaultImageIfEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
