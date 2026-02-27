package roleSvc

import "errors"

var (
	ErrRoleNotFound              = errors.New("场地权限角色不存在")
	ErrRoleNameRequired          = errors.New("角色名称不能为空")
	ErrRoleVenueIDInvalid        = errors.New("场地ID不合法")
	ErrRoleCreateInDB            = errors.New("数据库创建角色失败")
	ErrRoleUpdateInDB            = errors.New("数据库更新角色失败")
	ErrRoleQueryInDB             = errors.New("数据库查询角色失败")
	ErrRoleAccessGroupNotMatched = errors.New("权限组写入数量与请求不一致")
)
