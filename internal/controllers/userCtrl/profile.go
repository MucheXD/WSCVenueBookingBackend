package userCtrl

import (
	"errors"
	"regexp"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type UserProfileDTO struct {
	UID         string `json:"uid"`
	Username    string `json:"username"`
	SchoolID    string `json:"school_id"`
	RealName    string `json:"real_name"`
	PhoneNumber string `json:"phone_number"`
}

func GetSelfProfileHandler(c *gin.Context) {
	userID := c.GetString("UserID")
	if userID == "" {
		apiException.AbortWithException(c, apiException.Unauthorized)
		return
	}

	user, err := userSvc.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, userSvc.ErrUserNotFound) {
			apiException.AbortWithException(c, apiException.UserNotFound, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, UserProfileDTO{
		UID:         user.UID,
		Username:    user.Username,
		SchoolID:    user.SchoolID,
		RealName:    user.RealName,
		PhoneNumber: user.PhoneNumber,
	})
}

type AdminUserProfileDTO struct {
	UID          string `json:"uid"`
	Username     string `json:"username"`
	SchoolID     string `json:"school_id"`
	RealName     string `json:"real_name"`
	PhoneNumber  string `json:"phone_number"`
	RegisteredAt string `json:"registered_at"`
	UpdatedAt    string `json:"updated_at"`
	PermMap      uint64 `json:"perm_map"`
	PermVAGID    int    `json:"perm_vagid"`
}

func GetUserProfileHandler(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		apiException.AbortWithException(c, apiException.ParamError)
		return
	}

	user, err := userSvc.GetUserProfile(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, userSvc.ErrUserNotFound) {
			apiException.AbortWithException(c, apiException.UserNotFound, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, AdminUserProfileDTO{
		UID:          user.UID,
		Username:     user.Username,
		SchoolID:     user.SchoolID,
		RealName:     user.RealName,
		PhoneNumber:  user.PhoneNumber,
		RegisteredAt: formatRFC3339(user.RegisterTime),
		UpdatedAt:    formatRFC3339(user.UpdatedAt),
		PermMap:      user.PermMap,
		PermVAGID:    user.PermVAGID,
	})
}

type UserProfileUpdateDTO struct {
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
}

func UpdateSelfProfileHandler(c *gin.Context) {
	var req UserProfileUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	if !isValidPhoneNumber(req.PhoneNumber) {
		apiException.AbortWithException(c, apiException.PhoneNumberInvalid)
		return
	}

	err := userSvc.UpdateUser(c.Request.Context(), c.GetString("UserID"),
		models.User{
			PhoneNumber: req.PhoneNumber,
			Username:    req.Username})

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

func formatRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
