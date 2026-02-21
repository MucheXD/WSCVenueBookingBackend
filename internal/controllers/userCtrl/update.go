package userCtrl

import (
	"regexp"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type UserEditForm struct {
	PhoneNumber string `form:"phone_number" binding:"required"`
}

func UpdateUserProfileHandler(c *gin.Context) {
	var req UserEditForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if !isValidPhoneNumber(req.PhoneNumber) {
		apiException.AbortWithException(c, apiException.PhoneNumberInvalid)
		return
	}

	err := userSvc.UpdateUser(c.Request.Context(), c.GetString("UserID"),
		models.User{PhoneNumber: req.PhoneNumber})

	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.SetSuccessJsonResponse(c, nil)
}

var phoneNumberRegex = regexp.MustCompile(`^\d{11}$`)

func isValidPhoneNumber(phoneNumber string) bool {
	return phoneNumberRegex.MatchString(phoneNumber)
}
