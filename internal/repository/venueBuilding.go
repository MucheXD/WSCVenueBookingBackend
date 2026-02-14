package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type VenueBuildingEntity struct {
	BuildingID int    `gorm:"column:building_id;primaryKey"`
	Name       string `gorm:"column:building_name"`
	CampusID   int    `gorm:"column:location_campus_id"`
}

func (VenueBuildingEntity) TableName() string {
	return "venue_buildings"
}

func CreateNewVenueBuilding(modelB *models.VenueBuilding) error {
	var buildingEntity VenueBuildingEntity
	buildingEntity.fromDomain(modelB)
	txDB := database.DB.Create(&buildingEntity)
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

func GetVenueBuildingByID(buildingID int) (*models.VenueBuilding, error) {
	var buildingEntity VenueBuildingEntity
	txDB := database.DB.
		Where(&VenueBuildingEntity{BuildingID: buildingID}).
		Take(&buildingEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return buildingEntity.toDomain(), nil
}

func DeleteVenueBuildingByID(buildingID int) error {
	txDB := database.DB.
		Where(&VenueBuildingEntity{BuildingID: buildingID}).
		Delete(&VenueBuildingEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (b *VenueBuildingEntity) fromDomain(modelB *models.VenueBuilding) {
	b.BuildingID = modelB.ID
	b.Name = modelB.Name
	b.CampusID = modelB.CampusID
}

func (b *VenueBuildingEntity) toDomain() *models.VenueBuilding {
	return &models.VenueBuilding{
		ID:       b.BuildingID,
		Name:     b.Name,
		CampusID: b.CampusID,
	}
}
