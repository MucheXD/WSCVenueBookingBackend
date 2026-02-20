package repository

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
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
	CoverImageToken string         `gorm:"column:cover_image_token"`
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
		Distinct("venues.*").
		Joins("INNER JOIN venue_accesses va ON venues.venue_id = va.venue_id").
		Where("va.vagid = ?", opts.VAGID)

	// 权限筛选
	if len(opts.Permissions) > 0 {
		permConditions := database.DB.Where("1 = 0") // 初始化为false
		for _, perm := range opts.Permissions {
			// TODO: 此处避免使用硬编码文本
			switch perm {
			case venuePermission.Reserve:
				permConditions = permConditions.Or("va.allow_reserve = ?", true)
			case venuePermission.Approval:
				permConditions = permConditions.Or("va.allow_approval = ?", true)
			case venuePermission.Manage:
				permConditions = permConditions.Or("va.allow_manage = ?", true)
			case venuePermission.Edit:
				permConditions = permConditions.Or("va.allow_edit = ?", true)
			}
		}
		query = query.Where(permConditions)
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

// GetBuildingsByCampusesWithVenues 获取包含指定场地的楼区和校区信息
func GetBuildingsWithCampusesByVenueIDs(venueIDs []int) ([]*models.VenueBuilding, []*models.VenueCampus, error) {

	// SQL设计：使用Preload预加载关联的campus信息，避免N+1查询
	// SELECT b.*, c.* FROM venue_buildings b
	// INNER JOIN venue_campuses c ON b.location_campus_id = c.campus_id
	// WHERE b.building_id IN (SELECT DISTINCT location_building_id FROM venues WHERE venue_id IN (?))

	if len(venueIDs) == 0 {
		return []*models.VenueBuilding{}, []*models.VenueCampus{}, nil
	}

	// 获取场地对应的楼区ID
	var buildingIDs []int
	if err := database.DB.Model(&VenueEntity{}).
		Where("venue_id IN ?", venueIDs).
		Distinct("location_building_id").
		Pluck("location_building_id", &buildingIDs).Error; err != nil {
		return nil, nil, err
	}

	if len(buildingIDs) == 0 {
		return []*models.VenueBuilding{}, []*models.VenueCampus{}, nil
	}

	// 获取楼区信息
	var buildingEntities []VenueBuildingEntity
	if err := database.DB.Model(&VenueBuildingEntity{}).
		Where("building_id IN ?", buildingIDs).
		Find(&buildingEntities).Error; err != nil {
		return nil, nil, err
	}

	buildings := make([]*models.VenueBuilding, 0, len(buildingEntities))
	campusIDMap := make(map[int]struct{})
	for _, entity := range buildingEntities {
		buildings = append(buildings, entity.toDomain())
		campusIDMap[entity.CampusID] = struct{}{}
	}

	// 获取校区ID列表
	campusIDs := make([]int, 0, len(campusIDMap))
	for id := range campusIDMap {
		campusIDs = append(campusIDs, id)
	}

	// 获取校区信息
	var campusEntities []VenueCampusEntity
	if err := database.DB.Model(&VenueCampusEntity{}).
		Where("campus_id IN ?", campusIDs).
		Find(&campusEntities).Error; err != nil {
		return nil, nil, err
	}

	campuses := make([]*models.VenueCampus, 0, len(campusEntities))
	for _, entity := range campusEntities {
		campuses = append(campuses, entity.toDomain())
	}

	return buildings, campuses, nil
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
	v.CoverImageToken = modelV.CoverImageToken
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
		CoverImageToken: v.CoverImageToken,
		IsActive:        v.IsActive,
	}
}
