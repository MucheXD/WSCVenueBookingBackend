package middlewares

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/gin-gonic/gin"
)

// UnifiedErrorHandler 统一处理业务错误并输出标准响应。
func UnifiedErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() { // 如果已经有响应输出，直接返回
			return
		}

		if len(c.Errors) == 0 {
			utils.SetSuccessJsonResponse(c, nil)
			return
		}

		err := c.Errors.Last().Err
		var apiExc *apiException.Exception
		if !errors.As(err, &apiExc) || apiExc == nil {
			// TODO 这不是预期行为，需要记录日志
			apiExc = apiException.UnknownError
		}
		utils.SetErrorJsonResponse(c, apiExc.HTTPStatus, apiExc.Code, apiExc.Msg)

		// TODO 这里可以添加日志记录功能，记录 apiExc 及调用链上的其他错误
	}
}
