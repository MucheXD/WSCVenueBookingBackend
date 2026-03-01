package notificationSvc

import "errors"

var (
	ErrNotificationNotFound         = errors.New("公告不存在")
	ErrNotificationContentRequired     = errors.New("公告标题不能为空")
	ErrNotificationTitleRequired = errors.New("公告内容不能为空")
	ErrNotificationCreateInDB       = errors.New("数据库创建公告失败")
	ErrNotificationUpdateInDB       = errors.New("数据库更新公告失败")
	ErrNotificationDeleteInDB       = errors.New("数据库删除公告失败")
	ErrNotificationQueryInDB        = errors.New("数据库查询公告失败")
	ErrAttachmentQueryInDB   = errors.New("数据库查询附件失败")
)