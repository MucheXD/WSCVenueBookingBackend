package models

import "time"

type User struct {
	UID          string
	PasswordHash string
	PasswordSalt string
	RegisterTime time.Time
	Username     string
	SchoolID     string
	RealName     string
	PermMap      uint64
	PermVAGID    int
}
