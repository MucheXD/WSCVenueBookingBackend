package venuePermission

// VenuePermission 定义场地权限类型
type VenuePermission string

const (
	Reserve   VenuePermission = "reserve"   // 可预约
	Approval  VenuePermission = "approval"  // 可审批
	Edit      VenuePermission = "edit"      // 可编辑
	Manage    VenuePermission = "manage"    // 可管理
)
