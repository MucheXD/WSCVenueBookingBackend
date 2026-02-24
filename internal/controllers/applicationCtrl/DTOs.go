package applicationCtrl

type timetableDTO struct {
	Begin string `json:"begin" binding:"required"`
	End   string `json:"end" binding:"required"`
}

type attachmentDTO struct {
	Index     int    `json:"index"`
	FileToken string `json:"file_token" binding:"required"`
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
}

type commentDTO struct {
	ID          int             `json:"id"`
	Text        string          `json:"text"`
	Attachments []attachmentDTO `json:"attachments"`
}

type createApplicationForm struct {
	ApplicationType       string          `json:"application_type"`
	TimeRequest           []timetableDTO  `json:"time_request" binding:"required"`
	EstimatedParticipants *int            `json:"estimated_participants"`
	Description           string          `json:"description"`
	Attachments           []attachmentDTO `json:"attachments"`
	ActivityName          string          `json:"activity_name"`
	ActivityOrganizer     string          `json:"activity_organizer"`
	ActivityCoordinator   []any           `json:"activity_coordinator"`
}

type reviewApplicationForm struct {
	Decision       string      `json:"decision" binding:"required"`
	Comment        *commentDTO `json:"comment"`
	KnownConflicts []int       `json:"known_conflicts"`
}

type applicationResponseDTO struct {
	ID                    int             `json:"id"`
	ApplicationType       string          `json:"application_type"`
	ApplicationStatus     string          `json:"application_status"`
	TimeRequest           []timetableDTO  `json:"time_request"`
	EstimatedParticipants *int            `json:"estimated_participants,omitempty"`
	Description           string          `json:"description,omitempty"`
	Attachments           []attachmentDTO `json:"attachments,omitempty"`
	ActivityName          string          `json:"activity_name,omitempty"`
	ActivityOrganizer     string          `json:"activity_organizer,omitempty"`
	ActivityCoordinator   any             `json:"activity_coordinator,omitempty"`
	Comments              []commentDTO    `json:"comments"`
}
