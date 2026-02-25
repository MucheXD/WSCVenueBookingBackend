package notificationSvc

import (
	"context"
	"fmt"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"gorm.io/gorm"
)

func CreateNotification(ctx context.Context, notification models.Notification,senderUID string) (int, error) {
	if notification.Content == "" {
		return 0, ErrNotificationContentRequired
	}
	if notification.Title == "" {
		return 0, ErrNotificationTitleRequired
	}

	notification.SenderUID=senderUID
	
	var createdID int
	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		id, err := repository.CreateNotificationTx(tx, &notification)
		if err != nil {
			return err
		}
		createdID = id
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrNotificationCreateInDB, err)
	}

	return createdID, nil

}

func DeleteNotification(ctx context.Context, notificationID int) error {
	if _, err := repository.GetNotificationByID(notificationID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotificationNotFound
		}
		return fmt.Errorf("%w: %w", ErrNotificationQueryInDB, err)
	}

	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := repository.SoftDeleteNotificationsTx(tx, notificationID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotificationDeleteInDB, err)
	}
	return nil
}

func UpdateNotification(ctx context.Context, updates *models.Notification,notificationID int) error {
	existingNotification, err := repository.GetNotificationByID(notificationID)
	if err != nil {
		return err
	}

	// 应用更新（只更新非零值字段）
	utils.UpdateField(&existingNotification.Content, updates.Content)
	utils.UpdateField(&existingNotification.Title, updates.Title)
	utils.UpdateField(&existingNotification.RecevierUID, updates.RecevierUID)
	utils.UpdateField(&existingNotification.ReleaseTime, updates.ReleaseTime)
	utils.UpdateField(&existingNotification.Status, updates.Status)

	err = repository.UpdateNotification(existingNotification)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotificationUpdateInDB, err)
	}

	return nil
}

func ListNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	notifications, err := repository.ListNotifications(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotificationQueryInDB, err)
	}
	return notifications, nil
}

func ListAdminNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	notifications, err := repository.ListAdminNotifications(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNotificationQueryInDB, err)
	}
	return notifications, nil
}
