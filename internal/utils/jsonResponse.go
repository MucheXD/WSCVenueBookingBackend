package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func SetJsonResponse(c *gin.Context, httpStatus int, code int, msg string, data any) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func SetSuccessJsonResponse(c *gin.Context, data any) {
	SetJsonResponse(c, http.StatusOK, http.StatusOK, "success", data)
}

func SetErrorJsonResponse(c *gin.Context, httpStatus int, code int, msg string) {
	SetJsonResponse(c, httpStatus, code, msg, nil)
}
