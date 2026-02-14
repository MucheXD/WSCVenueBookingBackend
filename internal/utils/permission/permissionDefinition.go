package permission

var (
	AllVenueEdit           uint64 = 1 << 1  // 允许所有场地的 Edit 权限（无视场地权限表）
	AllVenueManage         uint64 = 1 << 2  // 允许所有场地的 Manage 权限（无视场地权限表）
	AllVenueApproval       uint64 = 1 << 3  // 允许所有场地的 Approval 权限（无视场地权限表）
	ChangeUserVenueAccess  uint64 = 1 << 4  // 允许修改用户的场地权限（因可自提权，等同具有所有场地权限）
	UserManagement         uint64 = 1 << 5  // 允许用户管理（用户信息增删改查）
	SendSystemAnnouncement uint64 = 1 << 6  // 允许发送系统通知
	SendUserNotification   uint64 = 1 << 7  // 允许发送用户通知（站内信）
	ChangeUserPermission   uint64 = 1 << 62 // 允许修改用户权限（因可自提权，等同最高权限）
	AllowAll               uint64 = 1 << 63 // 最高权限，允许所有操作
)
