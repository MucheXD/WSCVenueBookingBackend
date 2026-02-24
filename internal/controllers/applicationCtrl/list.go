package applicationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func ListVenueApplicationsHandler(c *gin.Context) {
	venueID, ok := parsePathInt(c, "venue_id")
	if !ok {
		return
	}
	vagid, sysPermMap, ok := getPermissionContext(c)
	if !ok {
		return
	}
	// cmt: 已将场地下申请单列表权限校验前置到 Controller 层，规则为 Reserve/AllVenueReservation。
	if !hasVenueReservePermission(vagid, sysPermMap, venueID) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	applications, err := applicationSvc.ListVenueApplications(c.Request.Context(), venueID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.SetSuccessJsonResponse(c, toApplicationResponseList(applications))
}

func ListMyApplicationsHandler(c *gin.Context) {
	userID := c.GetString("UserID")
	applications, err := applicationSvc.ListUserApplications(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.SetSuccessJsonResponse(c, toApplicationResponseList(applications))
}
