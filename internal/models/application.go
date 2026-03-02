package models

import (
	"encoding/json"
	"time"
)

const (
	ApplicationTypeNormal = "normal"

	ApplicationStatusRequested  = "req"
	ApplicationStatusRejected   = "rej"
	ApplicationStatusApproved   = "apv"
	ApplicationStatusMaintained = "maint"
)

type Application struct {
	ID                     int
	VenueID                int
	ApplicantUID           string
	ApplicationType        string
	ApplicationStatus      string
	TimeRequest            []ApplicationTimeRequest
	EstimatedParticipants  *int
	DescriptionText        string
	ActivityName           string
	ActivityOrganizer      string
	HasAttachments         bool
	Attachments            []Attachment
	Comments               []ApplicationComment
	ActivityCoordinatorRaw json.RawMessage
	CreatedAt              time.Time
	ApprovalAt             *time.Time
}

// ApplicationStats 申请单统计数据
type ApplicationStats struct {
	AllApplications   int64
	AllApproved       int64
	AllRejected       int64
	AllPending        int64
	Last7Applications int64
	Last7Approved     int64
	Last7Rejected     int64
}

type ApplicationTimeRequest struct {
	Begin time.Time
	End   time.Time
}

type ApplicationComment struct {
	ID             int
	ApplicationID  int
	CommenterUID   string
	CommentText    string
	HasAttachments bool
	CommentTime    time.Time
	Attachments    []Attachment
}

type ApplicationApproval struct {
	ApplicationID  int
	Decision       string
	Comment        *ApplicationComment
	KnownConflicts []int
}
