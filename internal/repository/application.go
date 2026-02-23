package repository

import (
	"sort"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type ApplicationEntity struct {
	ID                    int            `gorm:"column:id;primaryKey"`
	VenueID               int            `gorm:"column:venue_id"`
	ApplicantUID          string         `gorm:"column:applicant_uid"`
	ApplicationStatus     string         `gorm:"column:application_status"`
	DescriptionText       string         `gorm:"column:description_text"`
	EstimatedParticipants *int           `gorm:"column:estimated_participants"`
	HasAttachments        bool           `gorm:"column:has_attachments"`
	ActivityName          string         `gorm:"column:activity_name"`
	ActivityOrganizer     string         `gorm:"column:activity_organizer"`
	ActivityCoordinator   []byte         `gorm:"column:activity_coordinator"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ApplicationEntity) TableName() string {
	return "applications"
}

type ApplicationCommentEntity struct {
	ID             int            `gorm:"column:id;primaryKey"`
	ApplicationID  int            `gorm:"column:application_id"`
	CommenterUID   string         `gorm:"column:commenter_uid"`
	CommentText    string         `gorm:"column:comment_text"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
	HasAttachments bool           `gorm:"column:has_attachments"`
}

func (ApplicationCommentEntity) TableName() string {
	return "application_comments"
}

type ApplicationTimeslotEntity struct {
	ID            int            `gorm:"column:id;primaryKey"`
	VenueID       int            `gorm:"column:venue_id"`
	StartTime     time.Time      `gorm:"column:start_time"`
	EndTime       *time.Time     `gorm:"column:end_time"`
	ApplicationID *int           `gorm:"column:application_id"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ApplicationTimeslotEntity) TableName() string {
	return "venue_timeslots"
}

func CreateApplicationWithTx(tx *gorm.DB, application *models.Application) (int, error) {
	entity := ApplicationEntity{
		VenueID:               application.VenueID,
		ApplicantUID:          application.ApplicantUID,
		ApplicationStatus:     application.ApplicationStatus,
		DescriptionText:       application.DescriptionText,
		EstimatedParticipants: application.EstimatedParticipants,
		HasAttachments:        len(application.Attachments) > 0,
		ActivityName:          application.ActivityName,
		ActivityOrganizer:     application.ActivityOrganizer,
		ActivityCoordinator:   application.ActivityCoordinatorRaw,
	}
	if err := tx.Create(&entity).Error; err != nil {
		return 0, err
	}

	if len(application.TimeRequest) > 0 {
		timeslotEntities := make([]ApplicationTimeslotEntity, 0, len(application.TimeRequest))
		for _, slot := range application.TimeRequest {
			begin := slot.Begin
			end := slot.End
			appID := entity.ID
			timeslotEntities = append(timeslotEntities, ApplicationTimeslotEntity{
				VenueID:       application.VenueID,
				StartTime:     begin,
				EndTime:       &end,
				ApplicationID: &appID,
			})
		}
		if err := tx.Create(&timeslotEntities).Error; err != nil {
			return 0, err
		}
	}

	for idx, attachment := range application.Attachments {
		attachment.Index = idx
		if err := CreateAttachmentWithTx(tx, &attachment, AttachmentBizTypeApplication, entity.ID, idx); err != nil {
			return 0, err
		}
	}

	return entity.ID, nil
}

func CreateApplicationCommentWithTx(tx *gorm.DB, comment *models.ApplicationComment) (int, error) {
	entity := ApplicationCommentEntity{
		ApplicationID:  comment.ApplicationID,
		CommenterUID:   comment.CommenterUID,
		CommentText:    comment.CommentText,
		HasAttachments: len(comment.Attachments) > 0,
	}
	if err := tx.Create(&entity).Error; err != nil {
		return 0, err
	}
	for idx, attachment := range comment.Attachments {
		attachment.Index = idx
		if err := CreateAttachmentWithTx(tx, &attachment, AttachmentBizTypeApplicationComment, entity.ID, idx); err != nil {
			return 0, err
		}
	}
	return entity.ID, nil
}

func GetApplicationByID(applicationID int) (*models.Application, error) {
	applications, err := buildApplications(database.DB.Where("id = ?", applicationID))
	if err != nil {
		return nil, err
	}
	if len(applications) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &applications[0], nil
}

func ListApplicationsByVenueID(venueID int) ([]models.Application, error) {
	return buildApplications(database.DB.Where("venue_id = ?", venueID))
}

func ListApplicationsByApplicantUID(applicantUID string) ([]models.Application, error) {
	return buildApplications(database.DB.Where("applicant_uid = ?", applicantUID))
}

func SoftDeleteApplicationCascadeWithTx(tx *gorm.DB, applicationID int) error {
	if err := tx.Model(&ApplicationTimeslotEntity{}).
		Where("application_id = ?", applicationID).
		Delete(&ApplicationTimeslotEntity{}).Error; err != nil {
		return err
	}

	var comments []ApplicationCommentEntity
	if err := tx.Model(&ApplicationCommentEntity{}).
		Where("application_id = ?", applicationID).
		Find(&comments).Error; err != nil {
		return err
	}

	for _, comment := range comments {
		if err := SoftDeleteAttachmentsByBizWithTx(tx, AttachmentBizTypeApplicationComment, comment.ID); err != nil {
			return err
		}
	}

	if err := tx.Model(&ApplicationCommentEntity{}).
		Where("application_id = ?", applicationID).
		Delete(&ApplicationCommentEntity{}).Error; err != nil {
		return err
	}

	if err := SoftDeleteAttachmentsByBizWithTx(tx, AttachmentBizTypeApplication, applicationID); err != nil {
		return err
	}

	return tx.Model(&ApplicationEntity{}).
		Where("id = ?", applicationID).
		Delete(&ApplicationEntity{}).Error
}

func UpdateApplicationStatusWithTx(tx *gorm.DB, applicationID int, status string) error {
	return tx.Model(&ApplicationEntity{}).
		Where("id = ?", applicationID).
		Update("application_status", status).Error
}

func BatchRejectApplicationsWithTx(tx *gorm.DB, applicationIDs []int) error {
	if len(applicationIDs) == 0 {
		return nil
	}
	return tx.Model(&ApplicationEntity{}).
		Where("id IN ?", applicationIDs).
		Where("application_status = ?", models.ApplicationStatusRequested).
		Update("application_status", models.ApplicationStatusRejected).Error
}

func GetConflictingApplicationIDs(applicationID int, candidateStatuses []string) ([]int, error) {
	if len(candidateStatuses) == 0 {
		return []int{}, nil
	}

	var conflictIDs []int
	err := database.DB.
		Table("applications a1").
		Distinct("a2.id").
		Joins("JOIN venue_timeslots t1 ON t1.application_id = a1.id").
		Joins("JOIN venue_timeslots t2 ON t2.venue_id = a1.venue_id").
		Joins("JOIN applications a2 ON a2.id = t2.application_id").
		Where("a1.id = ?", applicationID).
		Where("a1.deleted_at IS NULL").
		Where("a2.deleted_at IS NULL").
		Where("t1.deleted_at IS NULL").
		Where("t2.deleted_at IS NULL").
		Where("a2.id <> a1.id").
		Where("a2.application_status IN ?", candidateStatuses).
		Where("t1.start_time < COALESCE(t2.end_time, '9999-12-31 23:59:59')").
		Where("COALESCE(t1.end_time, '9999-12-31 23:59:59') > t2.start_time").
		Pluck("a2.id", &conflictIDs).Error
	if err != nil {
		return nil, err
	}
	sort.Ints(conflictIDs)
	return conflictIDs, nil
}

func buildApplications(scope *gorm.DB) ([]models.Application, error) {
	var appEntities []ApplicationEntity
	if err := scope.
		Model(&ApplicationEntity{}).
		Order("id DESC").
		Find(&appEntities).Error; err != nil {
		return nil, err
	}
	if len(appEntities) == 0 {
		return []models.Application{}, nil
	}

	appIDs := make([]int, 0, len(appEntities))
	for _, app := range appEntities {
		appIDs = append(appIDs, app.ID)
	}

	var timeslotEntities []ApplicationTimeslotEntity
	if err := database.DB.
		Model(&ApplicationTimeslotEntity{}).
		Where("application_id IN ?", appIDs).
		Order("application_id ASC, start_time ASC").
		Find(&timeslotEntities).Error; err != nil {
		return nil, err
	}
	timeslotMap := make(map[int][]models.ApplicationTimeRequest)
	for _, timeslot := range timeslotEntities {
		if timeslot.ApplicationID == nil || timeslot.EndTime == nil {
			continue
		}
		appID := *timeslot.ApplicationID
		timeslotMap[appID] = append(timeslotMap[appID], models.ApplicationTimeRequest{
			Begin: timeslot.StartTime,
			End:   *timeslot.EndTime,
		})
	}

	var appAttachmentEntities []AttachmentEntity
	if err := database.DB.
		Model(&AttachmentEntity{}).
		Where("biz_type = ?", AttachmentBizTypeApplication).
		Where("biz_id IN ?", appIDs).
		Order("biz_id ASC, biz_index ASC").
		Find(&appAttachmentEntities).Error; err != nil {
		return nil, err
	}
	appAttachmentMap := make(map[int][]models.Attachment)
	for _, attachment := range appAttachmentEntities {
		appAttachmentMap[attachment.BizID] = append(appAttachmentMap[attachment.BizID], attachment.toDomain())
	}

	var commentEntities []ApplicationCommentEntity
	if err := database.DB.
		Model(&ApplicationCommentEntity{}).
		Where("application_id IN ?", appIDs).
		Order("application_id ASC, created_at ASC").
		Find(&commentEntities).Error; err != nil {
		return nil, err
	}

	commentIDs := make([]int, 0, len(commentEntities))
	for _, comment := range commentEntities {
		commentIDs = append(commentIDs, comment.ID)
	}

	commentAttachmentMap := make(map[int][]models.Attachment)
	if len(commentIDs) > 0 {
		var commentAttachmentEntities []AttachmentEntity
		if err := database.DB.
			Model(&AttachmentEntity{}).
			Where("biz_type = ?", AttachmentBizTypeApplicationComment).
			Where("biz_id IN ?", commentIDs).
			Order("biz_id ASC, biz_index ASC").
			Find(&commentAttachmentEntities).Error; err != nil {
			return nil, err
		}
		for _, attachment := range commentAttachmentEntities {
			commentAttachmentMap[attachment.BizID] = append(commentAttachmentMap[attachment.BizID], attachment.toDomain())
		}
	}

	commentsByApplication := make(map[int][]models.ApplicationComment)
	for _, comment := range commentEntities {
		commentsByApplication[comment.ApplicationID] = append(commentsByApplication[comment.ApplicationID], models.ApplicationComment{
			ID:             comment.ID,
			ApplicationID:  comment.ApplicationID,
			CommenterUID:   comment.CommenterUID,
			CommentText:    comment.CommentText,
			CommentTime:    comment.CreatedAt,
			HasAttachments: comment.HasAttachments,
			Attachments:    commentAttachmentMap[comment.ID],
		})
	}

	applications := make([]models.Application, 0, len(appEntities))
	for _, entity := range appEntities {
		applicationType := models.ApplicationTypeNormal
		applications = append(applications, models.Application{
			ID:                     entity.ID,
			VenueID:                entity.VenueID,
			ApplicantUID:           entity.ApplicantUID,
			ApplicationType:        applicationType,
			ApplicationStatus:      entity.ApplicationStatus,
			TimeRequest:            timeslotMap[entity.ID],
			EstimatedParticipants:  entity.EstimatedParticipants,
			DescriptionText:        entity.DescriptionText,
			ActivityName:           entity.ActivityName,
			ActivityOrganizer:      entity.ActivityOrganizer,
			HasAttachments:         entity.HasAttachments,
			Attachments:            appAttachmentMap[entity.ID],
			Comments:               commentsByApplication[entity.ID],
			ActivityCoordinatorRaw: entity.ActivityCoordinator,
		})
	}
	return applications, nil
}
