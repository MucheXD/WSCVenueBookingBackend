package middlewares

import (
	"strings"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/gin-gonic/gin"
)

// authMiddleware 身份鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiException.AbortWithException(c, apiException.Unauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			apiException.AbortWithException(c, apiException.AuthInvalid)
			return
		}
		tokenString := parts[1]
		tokenData, err := utils.GetTokenData(tokenString)
		if err != nil {
			apiException.AbortWithException(c, apiException.AuthInvalid)
			return
		}
		c.Set("user_id", tokenData.UserID)
		c.Next()
	}
}
