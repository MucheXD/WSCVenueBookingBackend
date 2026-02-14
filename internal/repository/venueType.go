package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type VenueTypeEntity struct {
	TypeID      int    `gorm:"column:venue_type_id;primaryKey"`
	Name        string `gorm:"column:name_text"`
	Description string `gorm:"column:description_text"`
}

func (VenueTypeEntity) TableName() string {
	return "venue_types"
}

func CreateNewVenueType(modelT *models.VenueType) error {
	var typeEntity VenueTypeEntity
	typeEntity.fromDomain(modelT)
	txDB := database.DB.Create(&typeEntity)
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

func GetVenueTypeByID(typeID int) (*models.VenueType, error) {
	var typeEntity VenueTypeEntity
	txDB := database.DB.
		Where(&VenueTypeEntity{TypeID: typeID}).
		Take(&typeEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return typeEntity.toDomain(), nil
}

func DeleteVenueTypeByID(typeID int) error {
	txDB := database.DB.
		Where(&VenueTypeEntity{TypeID: typeID}).
		Delete(&VenueTypeEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (t *VenueTypeEntity) fromDomain(modelT *models.VenueType) {
	t.TypeID = modelT.ID
	t.Name = modelT.Name
	t.Description = modelT.Description
}

func (t *VenueTypeEntity) toDomain() *models.VenueType {
	return &models.VenueType{
		ID:          t.TypeID,
		Name:        t.Name,
		Description: t.Description,
	}
}
