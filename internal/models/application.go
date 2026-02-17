package models

type Application struct {
	VenueID             int
	ApplicantUID        string
	ApplicationStatus   string
	DescriptionText     string
	ActivityName        string
	ActivityOrganizer   string
	Attachments         []Attachment
	Comments            []ApplicationComment
	ActivityCoordinator ApplicationCoordinator
}

type ApplicationCoordinator struct {
	UID          string
	Name         string
	ContactEmail string
	ContactPhone string
}

type ApplicationComment struct {
	CommenterUID string
	CommentText  string
	CommentTime  string
	Attachments  []Attachment
}
