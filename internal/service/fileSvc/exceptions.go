package fileSvc

import "errors"

var (
	ErrFileParamInvalid           = errors.New("文件参数不合法")
	ErrFileTokenInvalid           = errors.New("文件令牌不合法")
	ErrFileHashInvalid            = errors.New("文件哈希不合法")
	ErrFileSizeInvalid            = errors.New("文件大小不合法")
	ErrFileHashMismatch           = errors.New("文件哈希校验失败")
	ErrFileSizeMismatch           = errors.New("文件大小校验失败")
	ErrFileHashSizeConflict       = errors.New("文件哈希与大小不匹配")
	ErrFileNotFound               = errors.New("文件不存在")
	ErrFileStorageNotConfigured   = errors.New("文件存储未配置")
	ErrFileStorageTypeUnsupported = errors.New("不支持的文件存储类型")
	ErrFileQueryInDB              = errors.New("数据库查询文件失败")
	ErrFileCreateInDB             = errors.New("数据库创建文件记录失败")
	ErrFileUploadToStorage        = errors.New("文件上传存储失败")
	ErrFileDownloadFromStorage    = errors.New("文件读取存储失败")
)
