package notificationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/notificationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func HasUnreadNotificationHandler(c *gin.Context) {
	userID := c.GetString("UserID")
	hasUnread, err := notificationSvc.HasUnreadNotification(c.Request.Context(), userID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.SetSuccessJsonResponse(c, hasUnread)
}