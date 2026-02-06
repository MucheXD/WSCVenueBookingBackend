package middlewares

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/gin-gonic/gin"

)

// ErrHandler 中间件用于处理请求错误。
// 如果存在错误，将返回相应的 JSON 响应。
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 向下执行请求
		c.Next()

		// 如果存在错误，则处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if err != nil {
				var apiErr *apiException.Error

				// 尝试将错误转换为 apiException
				ok := errors.As(err, &apiErr)

				// 如果转换失败，则使用 ServerError
				if !ok {
					apiErr = apiException.ServerError
				}

				utils.JsonErrorResponse(c, apiErr.Code, apiErr.Msg)
				return
			}
		}
	}
}

// HandleNotFound 404处理
func HandleNotFound(c *gin.Context) {
	err := apiException.NotFound
	utils.JsonErrorResponse(c, err.Code, err.Msg)
}