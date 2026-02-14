package models

type Application struct {
	VenueID             int
	ApplicantUID        int
	ApplicationStatus   string
	DescriptionText     string
	ActivityName        string
	ActivityOrganizer   string
	Attachments         []Attachment
	Comments            []ApplicationComment
	ActivityCoordinator ApplicationCoordinator
}

type ApplicationCoordinator struct {
	UID          int
	Name         string
	ContactEmail string
	ContactPhone string
}

type ApplicationComment struct {
	CommenterUID int
	CommentText  string
	CommentTime  string
	Attachments  []Attachment
}
