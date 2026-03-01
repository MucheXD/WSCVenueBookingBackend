package notificationCtrl

import "time"

type attachmentDTO struct {
	Index     int    `json:"index"`
	FileToken string `json:"file_token" binding:"required"`
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
}

type createNotificationForm struct {
	ReceiverUID string          `json:"receiver_uid"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	Attachments []attachmentDTO `json:"attachments"`
	Status      int             `json:"status"`
	ReleaseTime time.Time       `json:"release_time"`
}

type updateNotificationForm struct {
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	Attachments []attachmentDTO `json:"attachments"`
	Status      int             `json:"status"`
	ReleaseTime time.Time       `json:"release_time"`
}

type notificationResponseDTO struct {
	NotificationID int             `json:"notification_id"`
	SenderUID      string          `json:"sender_uid"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	Attachments    []attachmentDTO `json:"attachments"`
	ReleaseTime    time.Time       `json:"release_time"`
}

type sentNotificationResponseDTO struct {
	NotificationID int             `json:"notification_id"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	Attachments    []attachmentDTO `json:"attachments"`
	Status         int             `json:"status"`
	ReleaseTime    time.Time       `json:"release_time"`
}
