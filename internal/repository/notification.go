package repository

import (
	"encoding/json"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NotificationContentEntity struct {
	ID          int             `gorm:"column:id;primaryKey"`
	Type        int             `gorm:"column:type"`
	Title       string          `gorm:"column:title"`
	Content     string          `gorm:"column:content"`
	SenderUID   string          `gorm:"column:sender_uid"`
	Payload     json.RawMessage `gorm:"column:payload"`
	Status      int             `gorm:"column:status"`
	ReleaseTime time.Time       `gorm:"column:release_time"`
}

type NotificationTargetEntity struct {
	ID             int    `gorm:"column:id;primaryKey"`
	NotificationID int    `gorm:"column:notification_id"`
	ReceiverUID    string `gorm:"column:receiver_uid"`
	IsRead         bool   `gorm:"column:is_read"`
	Type           int    `gorm:"column:type"`
}

func (NotificationContentEntity) TableName() string {
	return "notification_contents"
}

func (NotificationTargetEntity) TableName() string {
	return "notification_targets"
}

func CreateNotificationContentTx(tx *gorm.DB, notification *models.Notification) (int, error) {
	entity := NotificationContentEntity{
		Type:        notification.Type,
		Title:       notification.Title,
		Content:     notification.Content,
		Payload:     notification.Payload,
		SenderUID:   notification.SenderUID,
		Status:      notification.Status,
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
		if err := CreateAttachmentsTx(tx, AttachmentBizTypeNotification, entity.ID, attachments); err != nil {
			return 0, err
		}
	}

	return entity.ID, nil
}

func CreateNotificationTargetTx(tx *gorm.DB, notification *models.Notification) error {
	if notification == nil {
		return gorm.ErrInvalidData
	}
	return CreateNotificationTargetsTx(tx, []models.Notification{*notification})
}

func CreateNotificationTargetsTx(tx *gorm.DB, notifications []models.Notification) error {
	if len(notifications) == 0 {
		return nil
	}

	entities := make([]NotificationTargetEntity, 0, len(notifications))
	for _, notification := range notifications {
		if notification.ID == 0 || notification.ReceiverUID == "" {
			continue
		}
		entities = append(entities, NotificationTargetEntity{
			NotificationID: notification.ID,
			ReceiverUID:    notification.ReceiverUID,
			IsRead:         false,
			Type:           notification.Type,
		})
	}

	if len(entities) == 0 {
		return nil
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "notification_id"}, {Name: "receiver_uid"}},
		DoNothing: true,
	}).Create(&entities).Error
}

func SoftDeleteNotificationsTx(tx *gorm.DB, notificationID int) error {
	if err := SoftDeleteBizAttachmentsTx(tx, AttachmentBizTypeNotification, []int{notificationID}); err != nil {
		return err
	}
	err := tx.Model(&NotificationTargetEntity{}).
		Where(&NotificationTargetEntity{NotificationID: notificationID}).
		Delete(&NotificationTargetEntity{}).Error
	if err != nil {
		return err
	}

	return tx.Model(&NotificationContentEntity{}).
		Where(&NotificationContentEntity{ID: notificationID}).
		Delete(&NotificationContentEntity{}).Error
}

func GetNotificationByID(notificationID int) (*models.Notification, error) {
	var appEntities []NotificationContentEntity

	if err := database.DB.Model(&NotificationContentEntity{}).
		Where(&NotificationContentEntity{ID: notificationID}).
		Find(&appEntities).Error; err != nil {
		return nil, err
	}
	if len(appEntities) == 0 {
		return &models.Notification{}, gorm.ErrRecordNotFound
	}
	notificationContent := appEntities[0]

	var appAttachmentEntities []AttachmentEntity
	if err := database.DB.
		Model(&AttachmentEntity{}).
		Where("biz_type = ?", AttachmentBizTypeNotification).
		Where("biz_id = ?", notificationContent.ID).
		Find(&appAttachmentEntities).Error; err != nil {
		return nil, err
	}

	attachmentSlice := make([]models.Attachment, 0, len(appAttachmentEntities))
	for _, appAttachmentEntity := range appAttachmentEntities {
		attachmentSlice = append(attachmentSlice, appAttachmentEntity.toDomain())

	}

	notification := models.Notification{
		ID:          notificationContent.ID,
		SenderUID:   notificationContent.SenderUID,
		Title:       notificationContent.Title,
		Content:     notificationContent.Content,
		Status:      notificationContent.Status,
		ReleaseTime: notificationContent.ReleaseTime,
		Attachments: attachmentSlice,
	}

	return &notification, nil
}

func UpdateNotification(notification *models.Notification) error {
	entity := NotificationContentEntity{
		ID:          notification.ID,
		Title:       notification.Title,
		Content:     notification.Content,
		SenderUID:   notification.SenderUID,
		Status:      notification.Status,
		ReleaseTime: notification.ReleaseTime,
	}
	if err := database.DB.Model(&NotificationContentEntity{}).
		Where(&NotificationContentEntity{ID: notification.ID}).
		Updates(&entity).Error; err != nil {
		return err
	}
	return nil
}
func GetUnreadSystemNotificationsByUserID(tx *gorm.DB, userID string) ([]int, error) {
	var notificationIDs []int
	if err := tx.Model(&NotificationContentEntity{}).
		Where("type = ?", 1).
		Where("status = ?", 1).
		Where("NOT EXISTS (SELECT 1 FROM notification_targets WHERE notification_targets.notification_id = notification_contents.id AND notification_targets.receiver_uid = ?)", userID).
		Pluck("id", &notificationIDs).Error; err != nil {
		return notificationIDs, err
	}
	return notificationIDs, nil
}

func ListNotifications(tx *gorm.DB, userID string) ([]models.Notification, error) {
	var appEntities []NotificationTargetEntity

	if err := tx.Model(&NotificationTargetEntity{}).
		Where(&NotificationTargetEntity{ReceiverUID: userID}).
		Order("notification_id DESC").
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

	var NotificationContentEntities []NotificationContentEntity
	if err := tx.
		Model(&NotificationContentEntity{}).
		Where("id IN ?", appIDs).
		Find(&NotificationContentEntities).Error; err != nil {
		return nil, err
	}

	var appAttachmentEntities []AttachmentEntity
	if err := tx.
		Model(&AttachmentEntity{}).
		Where("biz_type = ?", AttachmentBizTypeNotification).
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
	for _, entity := range NotificationContentEntities {
		notifications = append(notifications, models.Notification{
			ID:          entity.ID,
			SenderUID:   entity.SenderUID,
			ReceiverUID: userID,
			Title:       entity.Title,
			Content:     entity.Content,
			Status:      entity.Status,
			ReleaseTime: entity.ReleaseTime,
			Attachments: appAttachmentMap[entity.ID],
		})
	}
	return notifications, nil
}

func MarkRead(tx *gorm.DB, UserID string) error {
	err := tx.Model(&NotificationTargetEntity{}).
		Where("receiver_uid = ?", UserID).
		Where("is_read = ?", false).
		Update("is_read", true).Error
	if err != nil {
		return err
	}
	return nil
}

func ListSentNotifications(userID string) ([]models.Notification, error) {
	var appEntities []NotificationContentEntity

	if err := database.DB.Model(&NotificationContentEntity{}).
		Where(&NotificationContentEntity{SenderUID: userID}).
		Where(&NotificationContentEntity{Type: 1}).
		Order("CASE WHEN status=3 THEN 0 WHEN status=2 THEN 1 WHEN status=1 THEN 2 ELSE 3 END ").
		Order("id DESC").
		Find(&appEntities).Error; err != nil {
		return nil, err
	}
	if len(appEntities) == 0 {
		return []models.Notification{}, nil
	}

	appIDs := make([]int, 0, len(appEntities))
	for _, app := range appEntities {
		appIDs = append(appIDs, app.ID)
	}

	var appAttachmentEntities []AttachmentEntity
	if err := database.DB.
		Model(&AttachmentEntity{}).
		Where("biz_type = ?", AttachmentBizTypeNotification).
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
			ID:          entity.ID,
			SenderUID:   entity.SenderUID,
			ReceiverUID: "",
			Title:       entity.Title,
			Content:     entity.Content,
			Status:      entity.Status,
			ReleaseTime: entity.ReleaseTime,
			Attachments: appAttachmentMap[entity.ID],
		})
	}
	return notifications, nil
}

func GetUnreadNotificationsNum(userID string) (int, error) {
	var count int64
	err := database.DB.Model(&NotificationTargetEntity{}).
		Where(map[string]any{
			"is_read":      false,
			"receiver_uid": userID,
		}).
		Count(&count).Error

	if err != nil {
		return -1, err
	}

	var notificationIDs []int
	if err := database.DB.Model(&NotificationContentEntity{}).
		Where("type = ?", 1).
		Where("status = ?", 1).
		Where("sender_uid != ?", userID).
		Where("NOT EXISTS (SELECT 1 FROM notification_targets WHERE notification_targets.notification_id = notification_contents.id AND notification_targets.receiver_uid = ?)", userID).
		Pluck("id", &notificationIDs).Error; err != nil {
		return -1, err
	}
	num := int(count) + len(notificationIDs)
	return num, nil

}
