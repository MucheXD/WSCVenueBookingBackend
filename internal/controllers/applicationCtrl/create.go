package applicationCtrl

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func CreateApplicationHandler(c *gin.Context) {
	venueID, ok := parsePathInt(c, "venue_id")
	if !ok {
		return
	}
	userID := c.GetString("UserID")
	vagid, sysPermMap, ok := getPermissionContext(c)
	if !ok {
		return
	}
	// cmt: 已将 Reserve/AllVenueReservation 权限校验前置到 Controller 层。
	if !hasVenueReservePermission(vagid, sysPermMap, venueID) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	var req createApplicationForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	// 目前仅支持普通申请单
	// 后续可能追加快速申请等其他类型，届时再调整此处逻辑
	if req.ApplicationType != models.ApplicationTypeNormal {
		apiException.AbortWithException(c, apiException.ParamError, errors.New("application_type must be normal"))
		return
	}

	application, err := toApplicationModel(req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	applicationID, err := applicationSvc.CreateApplication(c.Request.Context(), venueID, userID, application)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	utils.SetSuccessJsonResponse(c, map[string]int{"application_id": applicationID})
}
