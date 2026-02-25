package userCtrl

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type UserChangePwdForm struct {
	VerifyType string `json:"verify_type" binding:"required"`
	VerifyData string `json:"verify_data" binding:"required"`
	NewPwd     string `json:"new_password" binding:"required"`
	NewPwdSalt string `json:"new_salt" binding:"required"`
}

func UserChangePwdHandler(c *gin.Context) {
	var req UserChangePwdForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if !isValidPasswordHash(req.NewPwd) {
		apiException.AbortWithException(c, apiException.PasswordOrSaltInvalid)
		return
	}
	if req.VerifyType == "password" {
		err := userSvc.ChangePasswordByOld(c.Request.Context(),
			req.VerifyData, c.GetString("UserID"), req.NewPwd, req.NewPwdSalt)
		if err != nil {
			if errors.Is(err, userSvc.ErrUserNotVerified) {
				apiException.AbortWithException(c, apiException.ChangePwdFailed, err)
				return
			}
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
	}
	utils.SetSuccessJsonResponse(c, nil)
}
