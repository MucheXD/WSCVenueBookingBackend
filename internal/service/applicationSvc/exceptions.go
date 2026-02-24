package applicationSvc

import "errors"

var (
	ErrApplicationNotFound         = errors.New("申请单不存在")
	ErrApplicationStatusInvalid    = errors.New("申请单状态不合法")
	ErrApplicationPermissionDenied = errors.New("申请单权限不足")
	ErrApplicationNoTimeRequest    = errors.New("申请时间段不能为空")
	ErrApplicationTimeRangeInvalid = errors.New("申请时间段不合法")
	ErrApplicationDecisionInvalid  = errors.New("审批结论不合法")
	ErrApplicationKnownConflictErr = errors.New("审批冲突集不一致")
	ErrApplicationCreateInDB       = errors.New("数据库创建申请单失败")
	ErrApplicationQueryInDB        = errors.New("数据库查询申请单失败")
	ErrApplicationUpdateInDB       = errors.New("数据库更新申请单失败")
	ErrApplicationDeleteInDB       = errors.New("数据库删除申请单失败")
)
