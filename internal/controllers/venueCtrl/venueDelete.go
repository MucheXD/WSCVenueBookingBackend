package venueCtrl

import (
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

// DeleteVenueHandler 删除场地（软删除）
// DELETE /api/venue/:venue_id
// 权限要求：对应场地的Edit权限 OR AllVenueEdit系统权限
func DeleteVenueHandler(c *gin.Context) {
	// 获取场地ID
	venueIDStr := c.Param("venue_id")
	venueID, err := strconv.Atoi(venueIDStr)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	// 权限检查：需要场地Edit权限或AllVenueEdit系统权限
	if !checkVenueEditPermission(c, venueID) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	// 调用服务层删除场地
	err = venueSvc.DeleteVenue(c.Request.Context(), venueID)
	if err != nil {
		if errors.Is(err, venueSvc.ErrVenueNotFound) {
			apiException.AbortWithException(c, apiException.NotFound, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, map[string]string{"status": "deleted"})
}
