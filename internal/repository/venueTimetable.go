package repository

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type VenueTimeslotEntity struct {
	VenueID       int        `gorm:"column:venue_id"`
	StartTime     time.Time  `gorm:"column:start_time"`
	EndTime       *time.Time `gorm:"column:end_time"`
	ApplicationID *int       `gorm:"column:application_id"`
	Status        string     `gorm:"column:status"`
}

func (VenueTimeslotEntity) TableName() string {
	return "venue_timeslots"
}

func CreateNewVenueTimetable(modelT *models.VenueTimetable) error {
	entities := timetableToTimeslotEntities(modelT)
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
		Model(&VenueTimeslotEntity{}).
		Where(&VenueTimeslotEntity{VenueID: venueID}).
		Order("start_time ASC").
		Find(&entities)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	if len(entities) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return timeslotEntitiesToTimetable(venueID, entities), nil
}

func DeleteVenueTimetableByVenueID(venueID int) error {
	txDB := database.DB.
		Model(&VenueTimeslotEntity{}).
		Where(&VenueTimeslotEntity{VenueID: venueID}).
		Delete(&VenueTimeslotEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion

// Relation between timetable and timeslot / model and entity:
// model-timetable = venueID + []timeslot
// entity = venueID + timeslot
// model-timetable = entity*N

func timetableToTimeslotEntities(modelT *models.VenueTimetable) []VenueTimeslotEntity {
	if modelT == nil || len(modelT.Timeslots) == 0 {
		return nil
	}
	entities := make([]VenueTimeslotEntity, 0, len(modelT.Timeslots))
	for _, slot := range modelT.Timeslots {
		var entity VenueTimeslotEntity
		entity.fromDomain(modelT.VenueID, slot)
		entities = append(entities, entity)
	}
	return entities
}

func timeslotEntitiesToTimetable(venueID int, entities []VenueTimeslotEntity) *models.VenueTimetable {
	if len(entities) == 0 {
		return nil
	}
	modelT := &models.VenueTimetable{
		VenueID:   venueID,
		Timeslots: make([]models.VenueTimeslot, 0, len(entities)),
	}
	for _, entity := range entities {
		slot := entity.toDomain()
		modelT.Timeslots = append(modelT.Timeslots, *slot)
	}
	return modelT
}

func (v *VenueTimeslotEntity) fromDomain(venueID int, timeslot models.VenueTimeslot) {
	v.VenueID = venueID
	v.StartTime = timeslot.StartTime
	v.Status = timeslot.Status
	if timeslot.ApplicationID != 0 {
		associatedApplicationID := timeslot.ApplicationID
		v.ApplicationID = &associatedApplicationID
	}
	if !timeslot.EndTime.IsZero() {
		v.EndTime = &timeslot.EndTime
	}
}

func (v *VenueTimeslotEntity) toDomain() *models.VenueTimeslot {
	modelT := &models.VenueTimeslot{
		StartTime:     v.StartTime,
		EndTime:       time.Time{},
		ApplicationID: 0,
		Status:        v.Status,
	}
	if v.ApplicationID != nil {
		modelT.ApplicationID = *v.ApplicationID
	}
	if v.EndTime != nil {
		modelT.EndTime = *v.EndTime
	}
	return modelT
}
