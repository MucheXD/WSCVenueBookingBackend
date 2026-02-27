package apiException

import "net/http"

var (
	UnknownError            = NewException(http.StatusInternalServerError, 200000, "未知错误")
	ServerError             = NewException(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError              = NewException(http.StatusBadRequest, 200501, "参数错误")
	NotFound                = NewException(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unauthorized            = NewException(http.StatusUnauthorized, 200401, "请先登录")
	AuthInvalid             = NewException(http.StatusUnauthorized, 200402, "登录异常，请重新登录")
	LoginTimeout            = NewException(http.StatusUnauthorized, 200403, "登录过程超时，请重试")
	LoginInvalid            = NewException(http.StatusUnauthorized, 200404, "登录信息无效，请重新登录")
	LoginFailed             = NewException(http.StatusUnauthorized, 201001, "登录失败，用户名或密码错误")
	UsernameAlreadyExists   = NewException(http.StatusBadRequest, 201002, "用户名已存在")
	UsernameInvalid         = NewException(http.StatusBadRequest, 201003, "用户名不合法")
	SchoolIDInvalid         = NewException(http.StatusBadRequest, 201004, "学校ID不合法")
	PasswordOrSaltInvalid   = NewException(http.StatusBadRequest, 201005, "密码或盐值不合法")
	PhoneNumberInvalid      = NewException(http.StatusBadRequest, 201006, "联系电话不合法")
	ChangePwdFailed         = NewException(http.StatusUnauthorized, 201007, "旧密码错误")
	SysPermNotSatisfied     = NewException(http.StatusForbidden, 201008, "系统权限不足")
	VenuePermNotSatisfied   = NewException(http.StatusForbidden, 201009, "场地权限不足")
	UserNotFound            = NewException(http.StatusOK, 201010, "用户不存在")
	ApplicationNotFound     = NewException(http.StatusNotFound, 202001, "申请单不存在")
	ApplicationInapprovable = NewException(http.StatusBadRequest, 202002, "该申请单不可审批")
)
