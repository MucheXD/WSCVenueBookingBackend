package notificationCtrl


type attachmentDTO struct {
	Index     int    `json:"index"`
	FileToken string `json:"file_token" binding:"required"`
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
}


type createNotificationForm struct {
	RecevierUID string `json:"recevier_uid"`
	Title string `json:"title"`
	Content string `json:"content"`
	Attachments []attachmentDTO `json:"attachments"`
	Status int `json:"status"`
	ReleaseTime string `json:"release_time"`
}

type updateNotificationForm struct {
	RecevierUID string   `json:"recevier_uid"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Attachments []attachmentDTO `json:"attachments"`
	Status      int      `json:"status"`
	ReleaseTime string   `json:"release_time"`
}

type notificationResponseDTO struct {
	NotificationID int `json:"notification_id"`
	SenderUID string `json:"sender_uid"`
	Title string `json:"title"`
	Content string `json:"content"`
	Attachments []attachmentDTO `json:"attachments"`
	ReleaseTime string `json:"release_time"`
}

type adminNotificationResponseDTO struct {
	NotificationID int `json:"notification_id"`
	RecevierUID string `json:"recevier_uid"`
	Title string `json:"title"`
	Content string `json:"content"`
	Attachments []attachmentDTO `json:"attachments"`
	Status int `json:"status"`
	ReleaseTime string `json:"release_time"`
}

