package systemPermission

import "github.com/gin-gonic/gin"

type SystemPermission uint64

var (
	BasicUser              SystemPermission = 1 << 0  // 基础权限，允许访问系统和查看个人信息
	CreateNewVenue         SystemPermission = 1 << 1  // 允许创建新的场地
	AllVenueReserve        SystemPermission = 1 << 2  // 允许所有场地的 Reservation 权限（无视场地权限表）
	AllVenueApproval       SystemPermission = 1 << 3  // 允许所有场地的 Approval 权限（无视场地权限表）
	AllVenueManage         SystemPermission = 1 << 4  // 允许所有场地的 Manage 权限（无视场地权限表）
	AllVenueEdit           SystemPermission = 1 << 5  // 允许所有场地的 Edit 权限（无视场地权限表）
	ChangeUserVenueAccess  SystemPermission = 1 << 6  // 允许修改用户的场地权限（因可自提权，等同具有所有场地权限）
	UserManagement         SystemPermission = 1 << 7  // 允许用户管理（用户信息增删改查）
	SendSystemAnnouncement SystemPermission = 1 << 8  // 允许发送系统通知
	SendUserNotification   SystemPermission = 1 << 9  // 允许发送用户通知（站内信）
	ChangeUserPermission   SystemPermission = 1 << 62 // 允许修改用户权限（因可自提权，等同最高权限）
	AllowAll               SystemPermission = 1 << 63 // 最高权限，允许所有操作
)

var (
	RegisterDefault = BasicUser
	PermTypeUser    = BasicUser
	PermTypeAdmin   = BasicUser |
		AllVenueReserve |
		AllVenueApproval |
		AllVenueManage |
		AllVenueEdit |
		ChangeUserVenueAccess |
		UserManagement |
		SendSystemAnnouncement |
		SendUserNotification
	PermTypeOperator      = AllowAll
	SysNoSpecialVenuePerm = 0
)

var SystemPermissionDisplayList = []gin.H{
	{
		"display_name": "基本权限",
		"detail":       "具有此权限的用户可以登录系统，查看和修改自己的个人信息等",
		"bit":          0,
	},
	{
		"display_name": "新建场地",
		"detail":       "具有此权限的用户可以创建新的场地",
		"bit":          1,
	},
	{
		"display_name": "访问与申请所有场地",
		"detail":       "具有此权限的用户可以访问系统内的所有场地，并且对所有场地具有申请权限（无视场地权限表）",
		"bit":          2,
	},
	{
		"display_name": "审批所有场地下属申请单",
		"detail":       "具有此权限的用户可以审批系统内所有场地的申请单（无视场地权限表）",
		"bit":          3,
	},
	{
		"display_name": "管理所有场地",
		"detail":       "具有此权限的用户可以管理系统内的所有场地，包括发布场地公告、创建维护等（无视场地权限表）",
		"bit":          4,
	},
	{
		"display_name": "编辑所有场地",
		"detail":       "具有此权限的用户可以编辑系统内的所有场地，包括更改标题、关键属性、删除等（无视场地权限表）",
		"bit":          5,
	},
	{
		"display_name": "修改用户的场地权限",
		"detail":       "具有此权限的用户可以修改系统内所有用户的场地权限（因可自提权，等同具有所有场地权限）",
		"bit":          6,
	},
	{
		"display_name": "用户管理",
		"detail":       "具有此权限的用户可以进行用户管理操作，包括用户信息的增删改查",
		"bit":          7,
	},
	{
		"display_name": "发送系统通知",
		"detail":       "具有此权限的用户可以向所有用户发送系统通知",
		"bit":          8,
	},
	{
		"display_name": "发送对单用户通知",
		"detail":       "具有此权限的用户可以向一个或多个用户发送站内信",
		"bit":          9,
	},
	{
		"display_name": "修改用户权限",
		"detail":       "具有此权限的用户可以修改系统内所有用户的权限（因可自提权，等同最高权限）",
		"bit":          62,
	},
}
