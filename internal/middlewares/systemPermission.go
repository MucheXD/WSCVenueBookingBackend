package middlewares

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/gin-gonic/gin"
)

func CheckSystemPermission(required ...systemPermission.SystemPermission) gin.HandlerFunc {
	return func(c *gin.Context) {
		permMapVal, exists := c.Get("SysPermissionMap")
		if !exists {
			apiException.AbortWithException(c,
				apiException.SysPermNotSatisfied, errors.New("用户无系统权限信息"))
			return
		}
		permMap, ok := permMapVal.(uint64)
		if !ok {
			apiException.AbortWithException(c,
				apiException.SysPermNotSatisfied, errors.New("用户系统权限信息格式错误"))
			return
		}
		if !systemPermission.Satisfy(permMap, required...) {
			apiException.AbortWithException(c,
				apiException.SysPermNotSatisfied)
			return
		}
	}
}
