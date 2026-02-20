package venuePermission

// VenuePerm 定义场地权限类型
type VenuePerm string

const (
	Reserve  VenuePerm = "reserve"  // 可预约
	Approval VenuePerm = "approval" // 可审批
	Edit     VenuePerm = "edit"     // 可编辑
	Manage   VenuePerm = "manage"   // 可管理
)
