package controllers

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type UserChangePwdForm struct {
	Username string `form:"username" binding:"required"`
	VerifyType string `form:"verify_type" binding:"required"`
	VerifyData string `form:"verify_data" binding:"required"`
	NewPassword string `form:"new_password" binding:"required"`
}

func UserChangePwdHandler(c *gin.Context){
	var req UserChangePwdForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if !isValidPasswordHash(req.NewPassword)  {
		apiException.AbortWithException(c, apiException.PasswordOrSaltInvalid)
		return
	}
	if(req.VerifyType=="password"){
		isValid,err:=userSvc.ChangePassword(c.Request.Context(),req.VerifyData,req.Username,req.NewPassword)
		if err!=nil{
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		if !isValid{
			apiException.AbortWithException(c, apiException.ChangePwdFailed, err)
			return
		}
	}
	
	utils.SetSuccessJsonResponse(c,nil)

}
