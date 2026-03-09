package userCtrl

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type userRegisterForm struct {
	Username     string `json:"username" binding:"required"`
	SchoolID     string `json:"school_id" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	PasswordSalt string `json:"password_salt" binding:"required"`
}

func UserRegisterHandler(c *gin.Context) {
	var req userRegisterForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	// 构造用户模型并注册
	user := &models.User{
		Username:     req.Username,
		SchoolID:     req.SchoolID,
		PasswordHash: req.PasswordHash,
		PasswordSalt: req.PasswordSalt,
	}
	createdUser, err := userSvc.RegisterUser(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, userSvc.ErrUsernameAlreadyExists) {
			apiException.AbortWithException(c, apiException.UsernameAlreadyExists, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.SetSuccessJsonResponse(c, map[string]string{"UID": createdUser.UID})
}
