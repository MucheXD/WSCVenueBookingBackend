package notificationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/notificationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)


func ListNotificationHandler(c *gin.Context) {
	userID:=c.GetString("UserID")
	notifications, err := notificationSvc.ListNotifications(c.Request.Context(),userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, toNotificationResponseList(notifications))
}
	
	

func ListAdminNotificationsHandler(c *gin.Context) {
	userID := c.GetString("UserID")
	notifications, err := notificationSvc.ListAdminNotifications(c.Request.Context(), userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.SetSuccessJsonResponse(c, toAdminNotificationResponseList(notifications))
}