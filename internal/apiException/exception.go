package apiException

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

)

type Error struct {
	Code  int
	Msg   string
}

var (
	ServerError               = NewError(200500, "系统异常，请稍后重试!")
	ParamError                = NewError(200501, "参数错误")
	NotFound                  = NewError(200404, http.StatusText(http.StatusNotFound))
)

func (e *Error) Error() string {
	return e.Msg
}

func NewError(Code int, msg string) *Error {
	return &Error{
		Code:  Code,
		Msg:   msg,
	}
}

// AbortWithException 用于返回自定义错误信息
func AbortWithException(c *gin.Context, apiError *Error, err error) {
	_ = c.AbortWithError(200, apiError)
}

// AbortWithError 用于兼容原始错误返回
func AbortWithError(c *gin.Context, err error) {
	var apiError *Error
	if errors.As(err, &apiError) {
		_ = c.AbortWithError(200, apiError)
	} else {
		_ = c.AbortWithError(200, ServerError)
	}
}
