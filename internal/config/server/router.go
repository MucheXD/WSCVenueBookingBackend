package server

import (
	userControllers "github.com/MucheXD/WSCVenueBookingBackend/internal/controllers/user"
	venueControllers "github.com/MucheXD/WSCVenueBookingBackend/internal/controllers/venue"
	middlewares "github.com/MucheXD/WSCVenueBookingBackend/internal/middleware"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/gin-gonic/gin"
)

func initRouter() {
	GinEngine.Use(middlewares.UnifiedErrorHandler())
	GinEngine.GET("/test", func(c *gin.Context) { c.String(200, "success") })
	
	// 用户相关路由
	GinEngine.GET("/api/get-login-session-salt",
		userControllers.StartLoginSessionHandler)
	GinEngine.POST("/api/login",
		userControllers.PasswordLoginHandler)
	GinEngine.POST("/api/register",
		userControllers.UserRegisterHandler)
	GinEngine.PATCH("/api/user/edit-profile", middlewares.AuthMiddleware(),
		userControllers.UpdateUserProfileHandler)

	// 场地相关路由
	// 创建场地 - 需要 CreateNewVenue 系统权限
	GinEngine.PUT("/api/venue",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.CreateNewVenue),
		venueControllers.CreateVenueHandler)

	// 更新场地 - 需要场地Edit权限或AllVenueEdit系统权限（在controller中动态检查）
	GinEngine.PUT("/api/venue/:venue_id",
		middlewares.AuthMiddleware(),
		venueControllers.UpdateVenueHandler)

	// 删除场地 - 需要场地Edit权限或AllVenueEdit系统权限（在controller中动态检查）
	GinEngine.DELETE("/api/venue/:venue_id",
		middlewares.AuthMiddleware(),
		venueControllers.DeleteVenueHandler)

	// 列出场地 - 需要登录即可（在controller中根据用户权限筛选）
	GinEngine.GET("/api/venue",
		middlewares.AuthMiddleware(),
		venueControllers.ListVenuesHandler)

	// 获取位置信息 - 需要登录即可
	GinEngine.GET("/api/venue/locations",
		middlewares.AuthMiddleware(),
		venueControllers.GetVenueLocationsHandler)
}
