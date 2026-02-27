package models

import (
	"encoding/json"
	"time"
)

type Venue struct {
	ID              int
	Name            string
	BuildingID      int
	TypeID          int
	Capacity        int
	Description     string
	CoverImageToken string
	EquipmentsRaw   json.RawMessage
	IsActive        bool
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
	AllowReserve  map[int]struct{}
	AllowApproval map[int]struct{}
	AllowEdit     map[int]struct{}
	AllowManage   map[int]struct{}
}

// HasReserve 检查是否拥有指定场地的预约权限
func (va *VenueAccess) HasReserve(venueID int) bool {
	if va == nil || va.AllowReserve == nil {
		return false
	}
	_, ok := va.AllowReserve[venueID]
	return ok
}

// HasApproval 检查是否拥有指定场地的审批权限
func (va *VenueAccess) HasApproval(venueID int) bool {
	if va == nil || va.AllowApproval == nil {
		return false
	}
	_, ok := va.AllowApproval[venueID]
	return ok
}

// HasEdit 检查是否拥有指定场地的编辑权限
func (va *VenueAccess) HasEdit(venueID int) bool {
	if va == nil || va.AllowEdit == nil {
		return false
	}
	_, ok := va.AllowEdit[venueID]
	return ok
}

// HasManage 检查是否拥有指定场地的管理权限
func (va *VenueAccess) HasManage(venueID int) bool {
	if va == nil || va.AllowManage == nil {
		return false
	}
	_, ok := va.AllowManage[venueID]
	return ok
}

type VenueTimeslot struct {
	StartTime     time.Time
	EndTime       time.Time
	ApplicationID int
}

type VenueTimetable struct {
	VenueID   int
	Timeslots []VenueTimeslot
}
