package venueCtrl

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"github.com/gin-gonic/gin"
)

// VenueDetailDTO 场地详情DTO
type VenueDetailDTO struct {
	VenueID         int             `json:"venue_id"`
	Name            string          `json:"name"`
	BuildingID      int             `json:"building_id"`
	TypeID          int             `json:"type_id"`
	DescriptionText string          `json:"description_text"`
	CoverImageToken string          `json:"cover_image_token"`
	Capacity        int             `json:"capacity"`
	Equipments      any             `json:"equipments"`
	Permissions     []string        `json:"permissions"`
	Pending         *int            `json:"pending_applications,omitempty"`
	Processed       *int            `json:"processed_applications,omitempty"`
	Attachments     []AttachmentDTO `json:"attachments"`
	Timetable       []TimeslotDTO   `json:"timetable"`
}

// AttachmentDTO 附件DTO
type AttachmentDTO struct {
	Type      string `json:"type"`
	Index     int    `json:"index"`
	FileToken string `json:"file_token"`
}

// TimeslotDTO 时间段DTO
type TimeslotDTO struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// ListVenuesHandler 列出场地
// GET /api/venue
func ListVenuesHandler(c *gin.Context) {
	// 从上下文获取用户的VAGID（由AuthMiddleware设置）
	vagidVal, exists := c.Get("VenueAccessGroupID")
	if !exists {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}
	vagid, ok := vagidVal.(int)
	if !ok {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	// 解析查询参数
	buildingIDs := queryArrayAtoi(c.QueryArray("b"))
	typeIDs := queryArrayAtoi(c.QueryArray("t"))
	permissionTags := c.Query("p")
	permissions := parsePermissionsFromString(permissionTags)
	filterByPendingApplications := parseBoolFromQuery(c.Query("a"))
	if hasPermissionFilter(permissionTags, 'R') {
		filterByPendingApplications = false
	}
	searchQuery := c.Query("s")
	offset, _ := strconv.Atoi(c.Query("o"))
	limit, _ := strconv.Atoi(c.Query("n"))

	// 设置默认分页大小
	if limit == 0 {
		limit = 12
	}

	// 构造查询选项
	// 提取系统权限字段
	permMapVal, exists := c.Get("SysPermissionMap")
	if !exists {
		apiException.AbortWithException(c, apiException.SysPermNotSatisfied)
		return
	}
	var sysPerm uint64
	if permMap, ok := permMapVal.(uint64); ok {
		sysPerm = permMap
	} else {
		apiException.AbortWithException(c, apiException.SysPermNotSatisfied)
		return
	}

	opts := venueSvc.VenueListOptions{
		BuildingIDs: buildingIDs,
		TypeIDs:     typeIDs,
		Permissions: permissions,
		SearchQuery: searchQuery,
		Offset:      offset,
		Limit:       limit,
		VAGID:       vagid,
		SysPerm:     sysPerm,
	}

	// 调用服务层查询场地列表
	applicationCounts := map[int]venueSvc.VenueApplicationCounts{}
	var err error
	var venues []*models.Venue
	if filterByPendingApplications {
		venues, applicationCounts, err = venueSvc.ListFullVenuesWithAppCount(c.Request.Context(), opts)
	} else {
		venues, err = venueSvc.ListFullVenues(c.Request.Context(), opts)
	}
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 构建DTO响应
	result := make([]VenueDetailDTO, 0, len(venues))
	now := time.Now()
	sevenDaysLater := now.AddDate(0, 0, 7)

	for _, venue := range venues {
		dto := VenueDetailDTO{
			VenueID:         venue.ID,
			Name:            venue.Name,
			BuildingID:      venue.BuildingID,
			TypeID:          venue.TypeID,
			DescriptionText: venue.Description,
			CoverImageToken: venue.CoverImageToken,
			Capacity:        venue.Capacity,
			Equipments:      unmarshalVenueEquipments(venue.EquipmentsRaw),
			Permissions:     getVenuePermissionStrings(vagid, venue.ID, sysPerm),
			Attachments:     []AttachmentDTO{},
			Timetable:       []TimeslotDTO{},
		}

		// 查询附件 (biz_type=1 代表venue)
		attachments, err := repository.GetAttachmentsByBiz(1, venue.ID)
		if err == nil {
			for idx, att := range attachments {
				dto.Attachments = append(dto.Attachments, AttachmentDTO{
					Type:      att.BizFileType,
					Index:     idx,
					FileToken: att.FileToken,
				})
			}
		}

		// 查询时间段（未来7天内，最多32条）
		timeslots, err := repository.GetVenueTimeslotsInRange(venue.ID, now, sevenDaysLater, 32)
		if err == nil {
			for _, ts := range timeslots {
				endTimeStr := ""
				if !ts.EndTime.IsZero() {
					endTimeStr = ts.EndTime.Format(time.RFC3339)
				}
				dto.Timetable = append(dto.Timetable, TimeslotDTO{
					Start: ts.StartTime.Format(time.RFC3339),
					End:   endTimeStr,
				})
			}
		}

		if filterByPendingApplications {
			pending := 0
			processed := 0
			if counts, ok := applicationCounts[venue.ID]; ok {
				pending = counts.PendingCount
				processed = counts.ProcessedCount
			}
			dto.Pending = &pending
			dto.Processed = &processed
		}

		result = append(result, dto)
	}

	utils.SetSuccessJsonResponse(c, result)
}

// queryArrayAtoi 从 QueryArray 输出的字符数组解析整数数组
func queryArrayAtoi(values []string) []int {
	if len(values) == 0 {
		return nil
	}
	result := make([]int, 0, len(values))
	for _, val := range values {
		if num, err := strconv.Atoi(val); err == nil {
			result = append(result, num)
		}
	}
	return result
}

// parsePermissionsFromString 解析权限字符串（如"RAE"）
func parsePermissionsFromString(str string) []venuePermission.VenuePerm {
	if str == "" {
		return nil
	}
	result := make([]venuePermission.VenuePerm, 0, len(str))
	for _, ch := range str {
		if ch == 'R' {
			result = append(result, venuePermission.Reserve)
		}
		if ch == 'A' {
			result = append(result, venuePermission.Approval)
		}
		if ch == 'M' {
			result = append(result, venuePermission.Manage)
		}
		if ch == 'E' {
			result = append(result, venuePermission.Edit)
		}
	}
	return result
}

func parseBoolFromQuery(value string) bool {
	if value == "" {
		return false
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return parsed
}

func hasPermissionFilter(raw string, required ...rune) bool {
	if len(required) == 0 {
		return false
	}
	if raw == "" {
		return false
	}
	allowed := make(map[rune]struct{}, len(required))
	for _, tag := range required {
		allowed[tag] = struct{}{}
	}
	hasMatchedTag := false
	for _, ch := range raw {
		if _, ok := allowed[ch]; !ok {
			return false
		}
		hasMatchedTag = true
	}
	return hasMatchedTag
}

// getVenuePermissionStrings 获取用户对特定场地的权限字符串列表
func getVenuePermissionStrings(vagid int, venueID int, sysPerm uint64) []string {
	permissions := []string{}
	if venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Reserve) ||
		systemPermission.Check(sysPerm, systemPermission.AllVenueReserve) {
		permissions = append(permissions, "Reserve")
	}
	if venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Approval) ||
		systemPermission.Check(sysPerm, systemPermission.AllVenueApproval) {
		permissions = append(permissions, "Approval")
	}
	if venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Manage) ||
		systemPermission.Check(sysPerm, systemPermission.AllVenueManage) {
		permissions = append(permissions, "Manage")
	}
	if venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Edit) ||
		systemPermission.Check(sysPerm, systemPermission.AllVenueEdit) {
		permissions = append(permissions, "Edit")
	}
	return permissions
}

func unmarshalVenueEquipments(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil
	}
	return result
}

type VenueSimpleDTO struct {
	VenueID    int      `json:"venue_id"`
	Name       string   `json:"name"`
	BuildingID int      `json:"building_id"`
	TypeID     int      `json:"type_id"`
	Permission []string `json:"permission"`
}

// ListVenueAccessBodiesHandler 列出可修改权限的场地（轻量数据）
// GET /api/venue/list
func ListVenueAccessBodiesHandler(c *gin.Context) {
	vagid, err := strconv.Atoi(c.Param("vagid"))
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	venues, err := venueSvc.ListVenueBodies(c.Request.Context())
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	result := make([]VenueSimpleDTO, 0, len(venues))
	for _, venue := range venues {
		result = append(result, VenueSimpleDTO{
			VenueID:    venue.ID,
			Name:       venue.Name,
			BuildingID: venue.BuildingID,
			TypeID:     venue.TypeID,
			Permission: getVenuePermissionStrings(vagid, venue.ID, uint64(systemPermission.SysNoSpecialVenuePerm)),
		})
	}

	utils.SetSuccessJsonResponse(c, result)
}
