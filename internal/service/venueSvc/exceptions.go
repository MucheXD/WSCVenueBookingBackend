package venueSvc

import "errors"

var (
	ErrVenueNotFound         = errors.New("场地不存在")
	ErrVenueNameRequired     = errors.New("场地名称不能为空")
	ErrVenueBuildingRequired = errors.New("场地楼区不能为空")
	ErrVenueBuildingInvalid  = errors.New("场地楼区不合法")
	ErrVenueEquipmentsInvalid = errors.New("场地设备JSON不合法")
	ErrVenueCreateInDB       = errors.New("数据库创建场地失败")
	ErrVenueUpdateInDB       = errors.New("数据库更新场地失败")
	ErrVenueDeleteInDB       = errors.New("数据库删除场地失败")
	ErrVenueQueryInDB        = errors.New("数据库查询场地失败")
	ErrAttachmentQueryInDB   = errors.New("数据库查询附件失败")
	ErrTimeslotQueryInDB     = errors.New("数据库查询时间段失败")
	ErrBuildingQueryInDB     = errors.New("数据库查询楼区失败")
	ErrCampusQueryInDB       = errors.New("数据库查询校区失败")
)
