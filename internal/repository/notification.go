package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type NotificationEntity struct {
	NotificationID          int    `gorm:"column:notification_id;primaryKey"`
	Title string    `gorm:"column:title"`
	Content string    `gorm:"column:content"`
	SenderUID string `gorm:"column:sender_uid"`
	RecevierUID     string    `gorm:"column:recevier_uid"`
	Status     int    `gorm:"column:status"`
	ReleaseTime      string    `gorm:"column:release_time"`
}

func (NotificationEntity) TableName() string {
	return "notifications"
}



func CreateNotificationTx(tx *gorm.DB, notification *models.Notification) (int, error) {
	entity := NotificationEntity{
		NotificationID: notification.ID,
		Title: notification.Title,
		Content: notification.Content,
		RecevierUID: notification.RecevierUID,
		SenderUID: notification.SenderUID,
		Status: notification.Status,
		ReleaseTime: notification.ReleaseTime,
	}
	if err := tx.Create(&entity).Error; err != nil {
		return 0, err
	}

	if len(notification.Attachments) > 0 {
		attachments := make([]models.Attachment, 0, len(notification.Attachments))
		for idx, attachment := range notification.Attachments {
			attachment.Index = idx
			attachments = append(attachments, attachment)
		}
		if err := CreateAttachmentsTx(tx, AttachmentBizTypeNotification, entity.NotificationID, attachments); err != nil {
			return 0, err
		}
	}

	return entity.NotificationID, nil
}

func SoftDeleteNotificationsTx(tx *gorm.DB, notificationID int) error {
	if err := SoftDeleteBizAttachmentsTx(tx, AttachmentBizTypeNotification, []int{notificationID}); err != nil {
		return err
	}

	return tx.Model(&NotificationEntity{}).
		Where(&NotificationEntity{NotificationID: notificationID}).
		Delete(&NotificationEntity{}).Error
}

func GetNotificationByID(notificationID int) (*models.Notification, error) {
	notifications, err := queryNotifications(database.DB.Where(&NotificationEntity{NotificationID: notificationID}))
	if err != nil {
		return nil, err
	}
	if len(notifications) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &notifications[0], nil
}

func UpdateNotification(notification *models.Notification) error {
	entity := NotificationEntity{
		NotificationID: notification.ID,
		Title: notification.Title,
		Content: notification.Content,
		RecevierUID: notification.RecevierUID,
		SenderUID: notification.SenderUID,
		Status: notification.Status,
		ReleaseTime: notification.ReleaseTime,
	}
	if err := database.DB.Model(&NotificationEntity{}).
		Where(&NotificationEntity{NotificationID: entity.NotificationID}).
		Updates(&entity).Error; err != nil {
		return err
	}
	return nil
}

func ListNotifications(userID string) ([]models.Notification, error) {
	return queryNotifications(database.DB.Where("status=?",1).Where(&NotificationEntity{RecevierUID:userID},).Or("recevier_uid IS NULL",))
}

func ListAdminNotifications(userID string) ([]models.Notification, error) {
	return queryNotifications(database.DB.Where(&NotificationEntity{SenderUID: userID}))
}

func queryNotifications(scope *gorm.DB) ([]models.Notification, error) {
	var appEntities []NotificationEntity

	if err := scope.
		Model(&NotificationEntity{}).
		Order("id DESC").
		Find(&appEntities).Error; err != nil {
		return nil, err
	}
	if len(appEntities) == 0 {
		return []models.Notification{}, nil
	}

	appIDs := make([]int, 0, len(appEntities))
	for _, app := range appEntities {
		appIDs = append(appIDs, app.NotificationID)
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

	notifications := make([]models.Notification, 0, len(appEntities))
	for _, entity := range appEntities {
		notifications = append(notifications, models.Notification{
			ID:                     entity.NotificationID,
			SenderUID:                entity.SenderUID,
			RecevierUID:           entity.RecevierUID,
			Title:        entity.Title,
			Content:      entity.Content,
			Status:       entity.Status,
			ReleaseTime:   entity.ReleaseTime,
			Attachments:            appAttachmentMap[entity.NotificationID],
		})
	}
	return notifications, nil
}