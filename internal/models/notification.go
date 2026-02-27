package models

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID          int
	Type        int
	Title       string
	Content     string
	SenderUID   string
	RecevierUID string
	Payload     json.RawMessage
	Attachments []Attachment
	Status      int
	ReleaseTime time.Time
}
