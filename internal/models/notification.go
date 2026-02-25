package models


type Notification struct{
	ID int 
	Title string
	Content string 
	SenderUID string 
	RecevierUID string 
	Attachments []Attachment 
	Status int 
	ReleaseTime string 
}
