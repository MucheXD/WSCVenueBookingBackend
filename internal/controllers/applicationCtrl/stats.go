package applicationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/gin-gonic/gin"
)

type statsAllDTO struct {
	Applications int64 `json:"applications"`
	Approved     int64 `json:"approved"`
	Rejected     int64 `json:"rejected"`
	Pending      int64 `json:"pending"`
}

type statsLastSevenDaysDTO struct {
	Applications int64 `json:"applications"`
	Approved     int64 `json:"approved"`
	Rejected     int64 `json:"rejected"`
}

type venueStatsResponseDTO struct {
	All           statsAllDTO           `json:"all"`
	LastSevenDays statsLastSevenDaysDTO `json:"last_seven_days"`
}

// GetVenueStatsHandler 获取场地申请单统计数据
// 权限：需具有 AllVenueManage 或 AllVenueEdit 系统权限
func GetVenueStatsHandler(c *gin.Context) {
	sysPermMapVal, exists := c.Get("SysPermissionMap")
	if !exists {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return
	}
	sysPermMap, ok := sysPermMapVal.(uint64)
	if !ok {
		apiException.AbortWithException(c, apiException.AuthInvalid)
		return
	}
	if !systemPermission.SatisfyAny(sysPermMap, systemPermission.AllVenueManage, systemPermission.AllVenueEdit) {
		apiException.AbortWithException(c, apiException.SysPermNotSatisfied)
		return
	}

	stats, err := applicationSvc.GetApplicationStats(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	utils.SetSuccessJsonResponse(c, venueStatsResponseDTO{
		All: statsAllDTO{
			Applications: stats.AllApplications,
			Approved:     stats.AllApproved,
			Rejected:     stats.AllRejected,
			Pending:      stats.AllPending,
		},
		LastSevenDays: statsLastSevenDaysDTO{
			Applications: stats.Last7Applications,
			Approved:     stats.Last7Approved,
			Rejected:     stats.Last7Rejected,
		},
	})
}
