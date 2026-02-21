package userCtrl

import (
	"errors"
	"regexp"

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
	// 校验前端传入的数据是否合法，避免不必要的数据库查询和潜在的安全风险
	// 这里的校验规则可以根据实际需求进行调整
	if !isValidUsername(req.Username) {
		apiException.AbortWithException(c, apiException.UsernameInvalid)
		return
	}
	if !isValidSchoolID(req.SchoolID) {
		apiException.AbortWithException(c, apiException.SchoolIDInvalid)
		return
	}
	if !isValidPasswordHash(req.PasswordHash) || !isValidPasswordSalt(req.PasswordSalt) {
		apiException.AbortWithException(c, apiException.PasswordOrSaltInvalid)
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

// 预编译正则表达式以提高性能
// ^[a-zA-Z0-9_]+$ 表示只允许字母、数字和下划线
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// IsValidUsername 校验用户名是否合法
// 规则：4-20个字符，仅限英文字母、数字和下划线，不含特殊符号
func isValidUsername(username string) bool {
	length := len(username)
	if length < 4 || length > 20 {
		return false
	}
	return usernameRegex.MatchString(username)
}

// 这里可以根据实际需求定义学校ID的格式规则
// 必须是数字，长度为12位
var schoolIDRegex = regexp.MustCompile(`^\d{12}$`)

func isValidSchoolID(schoolID string) bool {
	return schoolIDRegex.MatchString(schoolID)
}

var passwordHashRegex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
var passwordSaltRegex = regexp.MustCompile(`^[a-zA-Z0-9]{16}$`)

func isValidPasswordHash(passwordHash string) bool {
	return passwordHashRegex.MatchString(passwordHash)
}
func isValidPasswordSalt(passwordSalt string) bool {
	return passwordSaltRegex.MatchString(passwordSalt)
}
