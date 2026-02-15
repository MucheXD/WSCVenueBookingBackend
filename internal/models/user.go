package models

import "time"

type User struct {
	UID          int
	PasswordHash string
	PasswordSalt string
	CreateTime   time.Time
	Username     string
	SchoolID     string
	RealName     string
	PermMap      uint64
	PermVAGID    int
}
