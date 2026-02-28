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
	// 获取登录盐值 - 无需认证
	GinEngine.GET("/api/login-session-salt",
		userCtrl.StartLoginSessionHandler)
	// 用户登录 - 无需认证
	GinEngine.POST("/api/login",
		userCtrl.PasswordLoginHandler)
	// 用户注册 - 无需认证
	GinEngine.POST("/api/register",
		userCtrl.UserRegisterHandler)
	// 修改自身信息 - 需要登录
	GinEngine.PUT("/api/user/profile", middlewares.AuthMiddleware(),
		userCtrl.UpdateSelfProfileHandler)
	// 获取自身信息 - 需要登录
	GinEngine.GET("/api/user/profile", middlewares.AuthMiddleware(),
		userCtrl.GetSelfProfileHandler)
	// 获取用户信息 - 需要 UserManagement 系统权限
	GinEngine.GET("/api/user/profile/:uid", middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.UserManagement),
		userCtrl.GetUserProfileHandler)
	// 修改密码 - 需要登录
	GinEngine.POST("/api/user/change-password", middlewares.AuthMiddleware(),
		userCtrl.UserChangePwdHandler)
	// (批量)修改用户系统权限 - 需要 ChangeUserPermission 系统权限
	GinEngine.PUT("/api/user/system-permission", middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserPermission),
		userCtrl.UpdateUserSysPermHandler)
	// 获取系统权限列表 - 需要 ChangeUserPermission 系统权限
	GinEngine.GET("/api/system-permission",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserPermission),
		userCtrl.GetSystemPermissionListHandler)
	// (批量)修改用户场地权限组 - 需要 ChangeUserVenueAccess 系统权限
	GinEngine.PUT("/api/user/vag", middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserVenueAccess),
		userCtrl.UpdateUserVAGHandler)

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

	// 列出可修改权限的场地（轻量数据）- 需要 ChangeUserVenueAccess 系统权限
	GinEngine.GET("/api/role/:vagid/venue",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserVenueAccess),
		venueCtrl.ListVenueAccessBodiesHandler)

	// 场地权限角色组管理 - 需要 ChangeUserVenueAccess 系统权限
	GinEngine.GET("/api/role",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserVenueAccess),
		venueCtrl.ListRolesHandler)
	GinEngine.POST("/api/role",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserVenueAccess),
		venueCtrl.CreateRoleHandler)
	GinEngine.PUT("/api/role/:vagid",
		middlewares.AuthMiddleware(),
		middlewares.CheckSystemPermission(systemPermission.ChangeUserVenueAccess),
		venueCtrl.UpdateRoleHandler)

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
