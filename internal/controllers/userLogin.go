package controllers

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type passwordLoginForm struct {
	LoginName  string `json:"loginName" binding:"required"`
	LoginToken string `json:"loginToken" binding:"required"`
}

func PasswordLogin(c *gin.Context) {
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
	utils.SetSuccessJsonResponse(c, "Login successful")
}
