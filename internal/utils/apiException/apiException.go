package apiException

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Exception struct {
	Code       int
	Msg        string
	HTTPStatus int
}

func (e *Exception) Error() string {
	return e.Msg
}

var (
	UnknownError = NewException(http.StatusInternalServerError, 200000, "未知错误")
	ServerError  = NewException(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError   = NewException(http.StatusBadRequest, 200501, "参数错误")
	NotFound     = NewException(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unauthorized = NewException(http.StatusUnauthorized, 200401, "请先登录")
	AuthInvalid  = NewException(http.StatusUnauthorized, 200402, "登录异常，请重新登录")
)

func NewException(httpStatus int, code int, msg string) *Exception {
	return &Exception{
		Code:       code,
		Msg:        msg,
		HTTPStatus: httpStatus,
	}
}

// AbortWithException 用于返回自定义错误信息
func AbortWithException(c *gin.Context, apiException *Exception, traceErrs ...error) {
	if apiException == nil {
		apiException = ServerError
	}
	_ = c.Error(apiException)       // 首异常为包装后的业务异常
	for _, err := range traceErrs { // 记录调用链上的其他错误(若有)
		if err != nil {
			_ = c.Error(err)
		}
	}
	c.Abort()
}
