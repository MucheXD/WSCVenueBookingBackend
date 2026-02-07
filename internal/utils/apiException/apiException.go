package apiException

import (
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
