package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type VenueCampusEntity struct {
	CampusID int    `gorm:"column:campus_id;primaryKey"`
	Name     string `gorm:"column:campus_name"`
}

func (VenueCampusEntity) TableName() string {
	return "venue_campuses"
}

func CreateNewVenueCampus(modelC *models.VenueCampus) error {
	var campusEntity VenueCampusEntity
	campusEntity.fromDomain(modelC)
	txDB := database.DB.Create(&campusEntity)
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

func GetVenueCampusByID(campusID int) (*models.VenueCampus, error) {
	var campusEntity VenueCampusEntity
	txDB := database.DB.
		Where(&VenueCampusEntity{CampusID: campusID}).
		Take(&campusEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return campusEntity.toDomain(), nil
}

func DeleteVenueCampusByID(campusID int) error {
	txDB := database.DB.
		Where(&VenueCampusEntity{CampusID: campusID}).
		Delete(&VenueCampusEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (c *VenueCampusEntity) fromDomain(modelC *models.VenueCampus) {
	c.CampusID = modelC.ID
	c.Name = modelC.Name
}

func (c *VenueCampusEntity) toDomain() *models.VenueCampus {
	return &models.VenueCampus{
		ID:   c.CampusID,
		Name: c.Name,
	}
}
