package controllers

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type startLoginSessionForm struct {
	LoginName string `json:"login_name" binding:"required"`
}
type passwordLoginForm struct {
	LoginName  string `json:"login_name" binding:"required"`
	LoginToken string `json:"login_token" binding:"required"`
}

func StartLoginSessionHandler(c *gin.Context) {
	var req startLoginSessionForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	// 获取登录名对应的用户盐值
	userSalt, err := userSvc.GetUserSalt(req.LoginName)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 生成登录会话盐值
	sessionSalt, err := userSvc.GenerateLoginSessionSalt(c.Request.Context(), req.LoginName)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, gin.H{
		"user_salt":    userSalt,
		"session_salt": sessionSalt,
	})
}

func PasswordLoginHandler(c *gin.Context) {
	var req passwordLoginForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	isValid, err := userSvc.TryPasswordLogin(c.Request.Context(), req.LoginName, req.LoginToken)
	if err != nil {
		if errors.Is(err, userSvc.ErrLoginTokenExpired) {
			apiException.AbortWithException(c, apiException.LoginTimeout, err)
			return
		}
		if errors.Is(err, userSvc.ErrLoginTokenSaltSecretNotConfigured) ||
			errors.Is(err, userSvc.ErrCheckLoginSessionIDInRedis) ||
			errors.Is(err, userSvc.ErrStoreLoginSessionIDInRedis) {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		apiException.AbortWithException(c, apiException.LoginInvalid, err)
		return
	}
	if !isValid {
		apiException.AbortWithException(c, apiException.LoginFailed, nil)
		return
	}
	utils.SetSuccessJsonResponse(c, nil)
}
