package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type VenueAccessEntity struct {
	ID               int  `gorm:"column:id;primaryKey"`
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
	entities := fromDomain(modelA)
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
	return toDomain(entities), nil
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

// --- Entity-Domain Conversion ---

// 批量 domain -> entities
func fromDomain(modelA *models.VenueAccess) []VenueAccessEntity {
	entities := make([]VenueAccessEntity, 0)
	appendEntities := func(venues map[int]struct{}, setFunc func(*VenueAccessEntity)) {
		for venueID := range venues {
			entity := VenueAccessEntity{
				VAGID:   modelA.VAGID,
				VenueID: venueID,
			}
			setFunc(&entity)
			entities = append(entities, entity)
		}
	}
	appendEntities(modelA.AllowReservation, func(e *VenueAccessEntity) { e.AllowReservation = true })
	appendEntities(modelA.AllowApproval, func(e *VenueAccessEntity) { e.AllowApproval = true })
	appendEntities(modelA.AllowEdit, func(e *VenueAccessEntity) { e.AllowEdit = true })
	appendEntities(modelA.AllowManage, func(e *VenueAccessEntity) { e.AllowManage = true })
	return mergeVenueAccessEntities(entities)
}

// 批量 entities -> domain
func toDomain(entities []VenueAccessEntity) *models.VenueAccess {
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
