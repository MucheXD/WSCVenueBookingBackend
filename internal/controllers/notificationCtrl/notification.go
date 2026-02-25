package notificationCtrl

import "github.com/MucheXD/WSCVenueBookingBackend/internal/models"

// 传入申请单 -> 申请单模型
func toNotificationModel(req createNotificationForm) (models.Notification, error) {
	return models.Notification{
		Title: req.Title,
		Content: req.Content,
		RecevierUID: req.RecevierUID,
		ReleaseTime: req.ReleaseTime,
		Status: req.Status,
		Attachments:            toAttachmentModelList(req.Attachments),
	}, nil
}

func toUpdateNotificationModel(req updateNotificationForm) (models.Notification, error) {
	return models.Notification{
		Title: req.Title,
		Content: req.Content,
		RecevierUID: req.RecevierUID,
		ReleaseTime: req.ReleaseTime,
		Status: req.Status,
		Attachments:            toAttachmentModelList(req.Attachments),
	}, nil
}

// 申请单模型 -> 申请单传出
func toNotificationResponseList(notifications []models.Notification) []notificationResponseDTO {
	result := make([]notificationResponseDTO, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, notificationResponseDTO{
			NotificationID:        notification.ID,
			SenderUID: notification.SenderUID,
			Title: notification.Title,
			Content: notification.Content,
			ReleaseTime: notification.ReleaseTime,
			Attachments:           toAttachmentDTOList(notification.Attachments),

		})
	}
	return result
}

func toAdminNotificationResponseList(notifications []models.Notification) []adminNotificationResponseDTO {
	result := make([]adminNotificationResponseDTO, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, adminNotificationResponseDTO{
			NotificationID:        notification.ID,
			RecevierUID: notification.RecevierUID,
			Title: notification.Title,
			Content: notification.Content,
			Status: notification.Status,
			ReleaseTime: notification.ReleaseTime,
			Attachments:           toAttachmentDTOList(notification.Attachments),

		})
	}
	return result
}

// 传入附件表 -> 附件表模型
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

// 附件表模型 -> 附件表传出
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