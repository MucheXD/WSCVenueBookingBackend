package server

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/controllers"
	middlewares "github.com/MucheXD/WSCVenueBookingBackend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter() {
	GinEngine.Use(middlewares.UnifiedErrorHandler())
	GinEngine.GET("/test", func(c *gin.Context) { c.String(200, "success") })
	GinEngine.GET("/api/get-login-session-salt",
		controllers.StartLoginSessionHandler)
	GinEngine.POST("/api/login",
		controllers.PasswordLoginHandler)
	GinEngine.POST("/api/register",
		controllers.UserRegisterHandler)
}
