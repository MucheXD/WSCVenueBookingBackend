package notificationCtrl

import (
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/notificationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

func DeleteNotificationHandler(c *gin.Context) {
	notificationIDStr := c.Param("notification_id")
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	err = notificationSvc.DeleteNotification(c.Request.Context(),notificationID)
	if err != nil {
		if errors.Is(err, notificationSvc.ErrNotificationNotFound) {
			apiException.AbortWithException(c, apiException.NotFound, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, map[string]string{"status": "deleted"})
}