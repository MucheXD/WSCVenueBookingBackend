package models

import "time"

type User struct {
	UID          string
	PasswordHash string
	PasswordSalt string
	RegisterTime time.Time
	UpdatedAt    time.Time
	Username     string
	SchoolID     string
	PhoneNumber  string
	RealName     string
	PermMap      uint64
	PermVAGID    int
}
