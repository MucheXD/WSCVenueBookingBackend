package middlewares

import (
	"errors"
	"log/slog"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

// UnifiedErrorHandler 统一处理业务错误并输出标准响应。
func UnifiedErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先完成后续处理

		if c.Writer.Written() { // 如果已经有响应输出，直接返回
			return
		}

		if len(c.Errors) == 0 { // 没有错误，正常返回成功响应
			utils.SetSuccessJsonResponse(c, nil)
			return
		}

		err := c.Errors.Last().Err
		var apiExc *apiException.Exception
		if !errors.As(err, &apiExc) || apiExc == nil {
			slog.Warn("Unknown API Exception occurred",
				"Error", err,
				"Trace", c.Errors)
			apiExc = apiException.UnknownError
		}
		utils.SetErrorJsonResponse(c, apiExc.HTTPStatus, apiExc.Code, apiExc.Msg)
		slog.Info("API Exception occurred",
			"Exception", apiExc,
			"Trace", c.Errors)
	}
}
