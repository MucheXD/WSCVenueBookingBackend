package server

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/controllers/applicationCtrl"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/controllers/fileCtrl"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/controllers/userCtrl"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/controllers/venueCtrl"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/middlewares"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/gin-gonic/gin"
)

func initRouter() {
	GinEngine.Use(middlewares.UnifiedErrorHandler())
	GinEngine.GET("/test", func(c *gin.Context) { c.String(200, "success") })

	// 用户相关路由
	GinEngine.GET("/api/login-session-salt",
		userCtrl.StartLoginSessionHandler)
	GinEngine.POST("/api/login",
		userCtrl.PasswordLoginHandler)
	GinEngine.POST("/api/register",
		userCtrl.UserRegisterHandler)
	GinEngine.PUT("/api/user/profile", middlewares.AuthMiddleware(),
		userCtrl.UpdateUserProfileHandler)
	GinEngine.POST("/api/user/change-password", middlewares.AuthMiddleware(),
		userCtrl.UserChangePwdHandler)

	// 文件相关路由
	GinEngine.POST("/api/file",
		fileCtrl.UploadFileHandler)
	GinEngine.GET("/api/file/:filetoken",
		fileCtrl.DownloadFileHandler)

	// 场地相关路由
	// 创建场地 - 需要 CreateNewVenue 系统权限
	GinEngine.PUT("/api/venue",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.CreateNewVenue),
		venueCtrl.CreateVenueHandler)

	// 更新场地 - 需要场地Edit权限或AllVenueEdit系统权限（在controller中动态检查）
	GinEngine.PUT("/api/venue/:venue_id",
		middlewares.AuthMiddleware(),
		venueCtrl.UpdateVenueHandler)

	// 删除场地 - 需要场地Edit权限或AllVenueEdit系统权限（在controller中动态检查）
	GinEngine.DELETE("/api/venue/:venue_id",
		middlewares.AuthMiddleware(),
		venueCtrl.DeleteVenueHandler)

	// 列出场地 - 需要登录即可（在controller中根据用户权限筛选）
	GinEngine.GET("/api/venue",
		middlewares.AuthMiddleware(),
		venueCtrl.ListVenuesHandler)

	// 获取位置信息 - 需要登录即可
	GinEngine.GET("/api/venue/locations",
		middlewares.AuthMiddleware(),
		venueCtrl.GetVenueLocationsHandler)

	// 申请单相关路由
	GinEngine.POST("/api/venue/:venue_id/application",
		middlewares.AuthMiddleware(),
		applicationCtrl.CreateApplicationHandler)
	GinEngine.DELETE("/api/application/:application_id",
		middlewares.AuthMiddleware(),
		applicationCtrl.DeleteApplicationHandler)
	GinEngine.GET("/api/venue/:venue_id/application",
		middlewares.AuthMiddleware(),
		applicationCtrl.ListVenueApplicationsHandler)
	GinEngine.GET("/api/user/application",
		middlewares.AuthMiddleware(),
		applicationCtrl.ListMyApplicationsHandler)
	GinEngine.PUT("/api/application/:application_id",
		middlewares.AuthMiddleware(),
		applicationCtrl.ReviewApplicationHandler)
}
