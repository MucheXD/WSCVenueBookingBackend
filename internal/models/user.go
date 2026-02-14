package models

import "time"

type User struct {
	UID          int
	PasswordHash string
	CreateTime   time.Time
	Username     string
	SchoolID     string
	RealName     string
	PermType     string
	PermVAGID    int
}
