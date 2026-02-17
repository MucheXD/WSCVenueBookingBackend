package server

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/controllers"
	middlewares "github.com/MucheXD/WSCVenueBookingBackend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter() {
	GinEngine.GET("/test", func(c *gin.Context) { c.String(200, "success") })
	GinEngine.GET("/api/get-login-session-salt",
		middlewares.UnifiedErrorHandler(),
		controllers.StartLoginSessionHandler)
	GinEngine.POST("/api/login",
		middlewares.UnifiedErrorHandler(),
		controllers.PasswordLoginHandler)
}
