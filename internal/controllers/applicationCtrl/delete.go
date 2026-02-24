package applicationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func DeleteApplicationHandler(c *gin.Context) {
	applicationID, ok := parsePathInt(c, "application_id")
	if !ok {
		return
	}
	requesterUID := c.GetString("UserID")
	vagid, sysPermMap, ok := getPermissionContext(c)
	if !ok {
		return
	}

	application, err := applicationSvc.GetApplicationByID(c.Request.Context(), applicationID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	// cmt: 已将删除权限校验前置到 Controller 层，规则为 申请人本人 或 Manage/AllVenueManage。
	if !canDeleteApplication(requesterUID, vagid, sysPermMap, *application) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	if err := applicationSvc.DeleteApplication(c.Request.Context(), applicationID); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.SetSuccessJsonResponse(c, nil)
}
