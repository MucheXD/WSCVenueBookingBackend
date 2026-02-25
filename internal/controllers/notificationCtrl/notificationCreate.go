package notificationCtrl

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/notificationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)


func CreateNotificationHandler(c *gin.Context) {
	userID:=c.GetString("UserID")

	var req createNotificationForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	notification, err := toNotificationModel(req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	notificationID, err := notificationSvc.CreateNotification(c.Request.Context(),notification,userID)
	if err != nil {
		if errors.Is(err, notificationSvc.ErrNotificationTitleRequired) ||
			errors.Is(err, notificationSvc.ErrNotificationContentRequired){
			apiException.AbortWithException(c, apiException.ParamError, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, map[string]int{"notification_id": notificationID})
}