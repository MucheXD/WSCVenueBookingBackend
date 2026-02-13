// TODO: AI Generated Code UNCHECKED

package repository

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type VenueTimeslotEntity struct {
	ID        int        `gorm:"column:id;primaryKey"`
	VenueID   int        `gorm:"column:venue_id"`
	StartTime time.Time  `gorm:"column:start_time"`
	EndTime   *time.Time `gorm:"column:end_time"`
	Status    string     `gorm:"column:status"`
}

func (VenueTimeslotEntity) TableName() string {
	return "venue_timeslots"
}

func CreateNewVenueTimetable(modelT *models.VenueTimetable) error {
	entities := fromTimetableDomain(modelT)
	if len(entities) == 0 {
		return nil
	}
	if err := database.DB.Create(&entities).Error; err != nil {
		return err
	}
	return nil
}

func GetVenueTimetableByVenueID(venueID int) (*models.VenueTimetable, error) {
	var entities []VenueTimeslotEntity
	txDB := database.DB.
		Where(&VenueTimeslotEntity{VenueID: venueID}).
		Order("start_time ASC").
		Find(&entities)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	if len(entities) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return toTimetableDomain(venueID, entities), nil
}

func DeleteVenueTimetableByVenueID(venueID int) error {
	txDB := database.DB.
		Where(&VenueTimeslotEntity{VenueID: venueID}).
		Delete(&VenueTimeslotEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// --- Entity-Domain Conversion ---

func fromTimetableDomain(modelT *models.VenueTimetable) []VenueTimeslotEntity {
	if modelT == nil || len(modelT.TimeSlots) == 0 {
		return nil
	}
	entities := make([]VenueTimeslotEntity, 0, len(modelT.TimeSlots))
	for _, slot := range modelT.TimeSlots {
		entity := VenueTimeslotEntity{
			VenueID:   modelT.VenueID,
			StartTime: time.Unix(slot.StartTime, 0),
			Status:    slot.Status,
		}
		if slot.EndTime > 0 {
			endTime := time.Unix(slot.EndTime, 0)
			entity.EndTime = &endTime
		}
		entities = append(entities, entity)
	}
	return entities
}

func toTimetableDomain(venueID int, entities []VenueTimeslotEntity) *models.VenueTimetable {
	if len(entities) == 0 {
		return nil
	}
	modelT := &models.VenueTimetable{
		VenueID:   venueID,
		TimeSlots: make([]models.VenueTimeslot, 0, len(entities)),
	}
	for _, entity := range entities {
		slot := models.VenueTimeslot{
			StartTime: entity.StartTime.Unix(),
			Status:    entity.Status,
		}
		if entity.EndTime != nil {
			slot.EndTime = entity.EndTime.Unix()
		}
		modelT.TimeSlots = append(modelT.TimeSlots, slot)
	}
	return modelT
}
