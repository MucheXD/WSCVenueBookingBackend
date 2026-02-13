package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type VenueEntity struct {
	VenueID        int    `gorm:"column:venue_id;primaryKey"`
	Name           string `gorm:"column:name_text"`
	BuildingID     int    `gorm:"column:location_building_id"`
	TypeID         int    `gorm:"column:type_id"`
	Capacity       int    `gorm:"column:capacity"`
	Description    string `gorm:"column:description_text"`
	CoverImageFile string `gorm:"column:cover_image_file"`
	IsActive       bool   `gorm:"column:is_active"`
}

func (VenueEntity) TableName() string {
	return "venues"
}

func CreateNewVenue(modelV *models.Venue) error {
	var venueEntity VenueEntity
	venueEntity.fromDomain(modelV)
	if err := database.DB.Create(&venueEntity).Error; err != nil {
		return err
	}
	return nil
}

func GetVenueByID(venueID int) (*models.Venue, error) {
	var venueEntity VenueEntity
	txDB := database.DB.
		Where(&VenueEntity{VenueID: venueID}).
		Take(&venueEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return venueEntity.toDomain(), nil
}

func DeleteVenueByID(venueID int) error {
	txDB := database.DB.
		Where(&VenueEntity{VenueID: venueID}).
		Delete(&VenueEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (v *VenueEntity) fromDomain(modelV *models.Venue) {
	v.VenueID = modelV.ID
	v.Name = modelV.Name
	v.BuildingID = modelV.BuildingID
	v.TypeID = modelV.TypeID
	v.Capacity = modelV.Capacity
	v.Description = modelV.Description
	v.CoverImageFile = modelV.CoverImageFile
	v.IsActive = modelV.IsActive
}

func (v *VenueEntity) toDomain() *models.Venue {
	return &models.Venue{
		ID:             v.VenueID,
		Name:           v.Name,
		BuildingID:     v.BuildingID,
		TypeID:         v.TypeID,
		Capacity:       v.Capacity,
		Description:    v.Description,
		CoverImageFile: v.CoverImageFile,
		IsActive:       v.IsActive,
	}
}
