package applicationCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/applicationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
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

	if err := applicationSvc.DeleteApplication(c.Request.Context(), applicationID, requesterUID, vagid, sysPermMap); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.SetSuccessJsonResponse(c, nil)
}
