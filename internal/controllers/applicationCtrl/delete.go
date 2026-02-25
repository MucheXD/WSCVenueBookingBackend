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

	application, err := applicationSvc.GetApplicationBodyByID(c.Request.Context(), applicationID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if !canDeleteApplication(requesterUID, vagid, sysPermMap, application.ApplicantUID, application.VenueID) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	if err := applicationSvc.DeleteApplication(c.Request.Context(), applicationID); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.SetSuccessJsonResponse(c, nil)
}
