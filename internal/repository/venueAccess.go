package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type VenueAccessEntity struct {
	VAGID            int  `gorm:"column:vagid"`
	VenueID          int  `gorm:"column:venue_id"`
	AllowReservation bool `gorm:"column:allow_reservation"`
	AllowApproval    bool `gorm:"column:allow_approval"`
	AllowEdit        bool `gorm:"column:allow_edit"`
	AllowManage      bool `gorm:"column:allow_manage"`
}

func (VenueAccessEntity) TableName() string {
	return "venue_accesses"
}

func CreateNewVenueAccessGroup(modelA *models.VenueAccess) error {
	entities := convVenueAccessToEntity(modelA)
	if len(entities) == 0 {
		return nil
	}
	if err := database.DB.Create(&entities).Error; err != nil {
		return err
	}
	return nil
}

func GetVenueAccessGroupByID(vagID int) (*models.VenueAccess, error) {
	var entities []VenueAccessEntity
	txDB := database.DB.
		Where(&VenueAccessEntity{VAGID: vagID}).
		Find(&entities)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	if len(entities) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return convVenueAccessFromEntity(entities), nil
}

func DeleteVenueAccessGroupByID(vagID int) error {
	txDB := database.DB.
		Where(&VenueAccessEntity{VAGID: vagID}).
		Delete(&VenueAccessEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion

// model-VenueAccess = vagID + [venueID]*N*permissions
// entity-VenueAccessEntity = vagID + venueID + permissions
// model-VenueAccess = entity*N

func convVenueAccessToEntity(modelA *models.VenueAccess) []VenueAccessEntity {
	if modelA == nil {
		return nil
	}
	entities := make([]VenueAccessEntity, 0)
	// 定义内部函数以减少重复代码，setFunc 为设置权限字段的函数
	appendEntities := func(venues map[int]struct{}, setFunc func(*VenueAccessEntity)) {
		if len(venues) == 0 {
			return
		}
		for venueID := range venues { // 将具有该权限的 venueID 分离出来作为一条记录
			entity := VenueAccessEntity{
				VAGID:   modelA.VAGID,
				VenueID: venueID,
			}
			setFunc(&entity) // 设置对应的权限字段，由传入的函数完成
			entities = append(entities, entity)
		}
	}
	// 下述步骤分离每种权限的每个目标场地作为独立的一条实体记录
	// 结果条目数为各权限对应的目标场地数之和，每个实体只有一个 true 权限字段
	appendEntities(modelA.AllowReservation, func(e *VenueAccessEntity) { e.AllowReservation = true })
	appendEntities(modelA.AllowApproval, func(e *VenueAccessEntity) { e.AllowApproval = true })
	appendEntities(modelA.AllowEdit, func(e *VenueAccessEntity) { e.AllowEdit = true })
	appendEntities(modelA.AllowManage, func(e *VenueAccessEntity) { e.AllowManage = true })
	// 由于数据表设计为窄表，故需合并同属一个目标场地的多个权限的记录
	// 结果条目数为不同的目标场地数
	return mergeVenueAccessEntities(entities)
}

func convVenueAccessFromEntity(entities []VenueAccessEntity) *models.VenueAccess {
	if len(entities) == 0 {
		return nil
	}
	vagid := entities[0].VAGID
	modelA := &models.VenueAccess{
		VAGID:            vagid,
		AllowReservation: make(map[int]struct{}),
		AllowApproval:    make(map[int]struct{}),
		AllowEdit:        make(map[int]struct{}),
		AllowManage:      make(map[int]struct{}),
	}
	for _, entity := range entities {
		if entity.AllowReservation {
			modelA.AllowReservation[entity.VenueID] = struct{}{}
		}
		if entity.AllowApproval {
			modelA.AllowApproval[entity.VenueID] = struct{}{}
		}
		if entity.AllowEdit {
			modelA.AllowEdit[entity.VenueID] = struct{}{}
		}
		if entity.AllowManage {
			modelA.AllowManage[entity.VenueID] = struct{}{}
		}
	}
	return modelA
}

// 合并同 vagid 且 venueID 的权限（防止重复）
func mergeVenueAccessEntities(entities []VenueAccessEntity) []VenueAccessEntity {
	result := make(map[[2]int]*VenueAccessEntity) // key: [vagid, venueID]
	for _, e := range entities {
		key := [2]int{e.VAGID, e.VenueID}
		if exist, ok := result[key]; ok {
			exist.AllowReservation = exist.AllowReservation || e.AllowReservation
			exist.AllowApproval = exist.AllowApproval || e.AllowApproval
			exist.AllowEdit = exist.AllowEdit || e.AllowEdit
			exist.AllowManage = exist.AllowManage || e.AllowManage
		} else {
			copy := e
			result[key] = &copy
		}
	}
	merged := make([]VenueAccessEntity, 0, len(result))
	for _, v := range result {
		merged = append(merged, *v)
	}
	return merged
}
