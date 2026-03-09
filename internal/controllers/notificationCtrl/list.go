package notificationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/notificationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func ListNotificationHandler(c *gin.Context) {
	userID := c.GetString("UserID")
	notifications, err := notificationSvc.ListNotifications(c.Request.Context(), userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, toNotificationResponseList(notifications))
}

func ListSentNotificationsHandler(c *gin.Context) {
	userID := c.GetString("UserID")

	_, sysPermMap, ok := getPermissionContext(c)
	if !ok {
		return
	}

	if !hasNotificationPermission(sysPermMap) {
		apiException.AbortWithException(c, apiException.SysPermNotSatisfied)
		return
	}

	notifications, err := notificationSvc.ListSentNotifications(c.Request.Context(), userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.SetSuccessJsonResponse(c, toSentNotificationResponseList(notifications))
}
