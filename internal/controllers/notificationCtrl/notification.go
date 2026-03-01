package notificationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/gin-gonic/gin"
)

func toNotificationModel(req createNotificationForm) (models.Notification, error) {
	return models.Notification{
		Title:       req.Title,
		Content:     req.Content,
		ReleaseTime: req.ReleaseTime,
		ReceiverUID: req.ReceiverUID,
		Status:      req.Status,
		Attachments: toAttachmentModelList(req.Attachments),
	}, nil
}

func toUpdateNotificationModel(req updateNotificationForm) (models.Notification, error) {
	return models.Notification{
		Title:       req.Title,
		Content:     req.Content,
		ReleaseTime: req.ReleaseTime,
		Status:      req.Status,
		Attachments: toAttachmentModelList(req.Attachments),
	}, nil
}

func toNotificationResponseList(notifications []models.Notification) []notificationResponseDTO {
	result := make([]notificationResponseDTO, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, notificationResponseDTO{
			NotificationID: notification.ID,
			SenderUID:      notification.SenderUID,
			Title:          notification.Title,
			Content:        notification.Content,
			ReleaseTime:    notification.ReleaseTime,
			Attachments:    toAttachmentDTOList(notification.Attachments),
		})
	}
	return result
}

func toSentNotificationResponseList(notifications []models.Notification) []sentNotificationResponseDTO {
	result := make([]sentNotificationResponseDTO, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, sentNotificationResponseDTO{
			NotificationID: notification.ID,
			Title:          notification.Title,
			Content:        notification.Content,
			Status:         notification.Status,
			ReleaseTime:    notification.ReleaseTime,
			Attachments:    toAttachmentDTOList(notification.Attachments),
		})
	}
	return result
}

func toAttachmentModelList(values []attachmentDTO) []models.Attachment {
	if len(values) == 0 {
		return []models.Attachment{}
	}
	result := make([]models.Attachment, 0, len(values))
	for _, value := range values {
		result = append(result, models.Attachment{
			Index:       value.Index,
			FileToken:   value.FileToken,
			BizFileType: value.FileType,
			BizFileName: value.FileName,
		})
	}
	return result
}

func toAttachmentDTOList(values []models.Attachment) []attachmentDTO {
	if len(values) == 0 {
		return []attachmentDTO{}
	}
	result := make([]attachmentDTO, 0, len(values))
	for _, value := range values {
		result = append(result, attachmentDTO{
			Index:     value.Index,
			FileToken: value.FileToken,
			FileType:  value.BizFileType,
			FileName:  value.BizFileName,
		})
	}
	return result
}

func getPermissionContext(c *gin.Context) (int, uint64, bool) {
	vagidVal, exists := c.Get("VenueAccessGroupID")
	if !exists {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}
	vagid, ok := vagidVal.(int)
	if !ok {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}

	sysPermVal, exists := c.Get("SysPermissionMap")
	if !exists {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}
	sysPermMap, ok := sysPermVal.(uint64)
	if !ok {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return 0, 0, false
	}
	return vagid, sysPermMap, true
}

func hasNotificationPermission(sysPermMap uint64) bool {
	if systemPermission.Check(sysPermMap, systemPermission.SendSystemAnnouncement) {
		return true
	}
	if systemPermission.Check(sysPermMap, systemPermission.SendUserNotification) {
		return true
	}
	if systemPermission.Check(sysPermMap, systemPermission.ChangeUserPermission) {
		return true
	}
	return systemPermission.Check(sysPermMap, systemPermission.AllowAll)
}
