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
	ID                  int
	VenueID             int
	ApplicantUID        string
	ApplicationType     string
	ApplicationStatus   string
	TimeRequest         []ApplicationTimeRequest
	EstimatedParticipants *int
	DescriptionText     string
	ActivityName        string
	ActivityOrganizer   string
	HasAttachments      bool
	Attachments         []Attachment
	Comments            []ApplicationComment
	ActivityCoordinatorRaw json.RawMessage
}

type ApplicationTimeRequest struct {
	Begin time.Time
	End   time.Time
}

type ApplicationComment struct {
	ID          int
	ApplicationID int
	CommenterUID string
	CommentText  string
	HasAttachments bool
	CommentTime  time.Time
	Attachments  []Attachment
}

type ApprovalDecision string

const (
	ApprovalDecisionRejected ApprovalDecision = "rejected"
	ApprovalDecisionApproved ApprovalDecision = "approved"
)

type ApplicationApproval struct {
	ApplicationID   int
	Decision        ApprovalDecision
	Comment         *ApplicationComment
	KnownConflicts  []int
}
