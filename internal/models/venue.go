package models

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
	VAGID         int
	TargetVenueID int
	// Note: Using map[int]struct{} to represent sets of Venues for each permission type
	// 注意：出于索引效率使用 map[int]struct{} 请按 Array 理解
	AllowReservation map[int]struct{}
	AllowApproval    map[int]struct{}
	AllowEdit        map[int]struct{}
	AllowManage      map[int]struct{}
}

type VenueTimeslot struct {
	StartTime int64
	EndTime   int64
	Status    string
}

type VenueTimetable struct {
	VenueID   int
	TimeSlots []VenueTimeslot
}
