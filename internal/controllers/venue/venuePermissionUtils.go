package controllers

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"github.com/gin-gonic/gin"
)

// checkVenueEditPermission 检查用户是否有场地编辑权限
// 需要：场地Edit权限 OR AllVenueEdit系统权限
func checkVenueEditPermission(c *gin.Context, venueID int) bool {
	// 检查系统权限
	sysPermMapVal, exists := c.Get("SysPermissionMap")
	if exists {
		if sysPermMap, ok := sysPermMapVal.(uint64); ok {
			if systemPermission.Check(sysPermMap, systemPermission.AllVenueEdit) {
				return true
			}
		}
	}

	// 检查场地权限
	vagidVal, exists := c.Get("VenueAccessGroupID")
	if !exists {
		return false
	}
	vagid, ok := vagidVal.(int)
	if !ok {
		return false
	}

	return venuePermission.CheckVenuePermission(vagid, venueID, venuePermission.Edit)
}
