package middlewares

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"github.com/gin-gonic/gin"
)

// CheckVenuePermission 场地权限检查中间件
// getVenueID 是一个函数，用于从 gin.Context 中提取 venueID
// perms 是需要检查的权限列表，只要满足其中一个权限即可通过
func CheckVenuePermission(venueID int, perms ...venuePermission.VenuePerm) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 中获取 VenueAccessGroupID（由 AuthMiddleware 设置）
		vagidVal, exists := c.Get("VenueAccessGroupID")
		if !exists {
			apiException.AbortWithException(c,
				apiException.VenuePermNotSatisfied, errors.New("用户无场地权限信息"))
			return
		}
		vagid, ok := vagidVal.(int)
		if !ok {
			apiException.AbortWithException(c,
				apiException.VenuePermNotSatisfied, errors.New("用户场地权限信息格式错误"))
			return
		}
		// 检查是否满足任意一个权限
		hasPermission := false
		for _, perm := range perms {
			if venuePermission.CheckVenuePermission(vagid, venueID, perm) {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
			return
		}
		c.Next()
	}
}
