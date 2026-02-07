package apiException

import "net/http"

var (
	UnknownError = NewException(http.StatusInternalServerError, 200000, "未知错误")
	ServerError  = NewException(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError   = NewException(http.StatusBadRequest, 200501, "参数错误")
	NotFound     = NewException(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unauthorized = NewException(http.StatusUnauthorized, 200401, "请先登录")
	AuthInvalid  = NewException(http.StatusUnauthorized, 200402, "登录异常，请重新登录")
)
