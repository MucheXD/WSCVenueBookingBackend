package systemPermission

type SystemPermission uint64

var (
	BasicUser              SystemPermission = 1 << 0  // 基础权限，允许访问系统和查看个人信息
	CreateNewVenue         SystemPermission = 1 << 1  // 允许创建新的场地
	AllVenueReservation    SystemPermission = 1 << 2  // 允许所有场地的 Reservation 权限（无视场地权限表）
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
		AllVenueReservation |
		AllVenueApproval |
		AllVenueManage |
		AllVenueEdit |
		ChangeUserVenueAccess |
		UserManagement |
		SendSystemAnnouncement |
		SendUserNotification
	PermTypeOperator = AllowAll
)
