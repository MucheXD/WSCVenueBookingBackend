package models

import "time"

type Venue struct {
	ID             int
	Name           string
	BuildingID     int
	TypeID         int
	Capacity       int
	Description    string
	CoverImageFile string
	IsActive       bool
}

type VenueType struct {
	ID          int
	Name        string
	Description string
}

type VenueBuilding struct {
	ID       int
	Name     string
	CampusID int
}

type VenueCampus struct {
	ID   int
	Name string
}

type VenueAccess struct {
	VAGID int
	// Note: Using map[int]struct{} to represent sets of Venues for each permission type
	// 注意：出于索引效率使用 map[int]struct{} 请按 Array 理解
	AllowReservation map[int]struct{}
	AllowApproval    map[int]struct{}
	AllowEdit        map[int]struct{}
	AllowManage      map[int]struct{}
}

type VenueTimeslot struct {
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

type VenueTimetable struct {
	VenueID   int
	Timeslots []VenueTimeslot
}
