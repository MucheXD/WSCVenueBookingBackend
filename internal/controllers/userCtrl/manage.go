package userCtrl

import (
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/systemPermission"
	"github.com/gin-gonic/gin"
)

type userListItemDTO struct {
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

func ListUsersHandler(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("o", "0"))
	if err != nil || offset < 0 {
		apiException.AbortWithException(c, apiException.ParamError)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("n", "10"))
	if err != nil || limit <= 0 {
		apiException.AbortWithException(c, apiException.ParamError)
		return
	}

	users, err := userSvc.ListUsers(c.Request.Context(), offset, limit)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	result := make([]userListItemDTO, 0, len(users))
	for _, user := range users {
		result = append(result, userListItemDTO{
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

	utils.SetSuccessJsonResponse(c, result)
}

type UpdateUserSysPermDTO struct {
	UserID  []string `json:"uids" binding:"required"`
	SysPerm uint64   `json:"system_permission" binding:"required"`
}

type UpdateUserVAGDTO struct {
	UserID []string `json:"uids" binding:"required"`
	VAGID  int      `json:"vagid" binding:"required"`
}

// GetSystemPermissionListHandler 获取系统权限列表
// 此函数的内容写死，不从数据库读取，因为系统权限是写死在代码里的
func GetSystemPermissionListHandler(c *gin.Context) {
	utils.SetSuccessJsonResponse(c, systemPermission.SystemPermissionDisplayList)
}

func UpdateUserSysPermHandler(c *gin.Context) {
	var req UpdateUserSysPermDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	if len(req.UserID) == 0 {
		apiException.AbortWithException(c, apiException.ParamError)
		return
	}

	if len(req.UserID) == 0 {
		apiException.AbortWithException(c, apiException.ParamError)
		return
	}

	err := userSvc.BatchUpdateUsersSystemPermission(c.Request.Context(), req.UserID, uint64(req.SysPerm))
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, nil)
}

func UpdateUserVAGHandler(c *gin.Context) {
	var req UpdateUserVAGDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	if len(req.UserID) == 0 {
		apiException.AbortWithException(c, apiException.ParamError)
		return
	}

	err := userSvc.BatchUpdateUsersVenueAccessGroup(c.Request.Context(), req.UserID, req.VAGID)
	if err != nil {
		if errors.Is(err, userSvc.ErrVenueAccessGroupInvalid) || errors.Is(err, userSvc.ErrVenueAccessGroupNotFound) {
			apiException.AbortWithException(c, apiException.ParamError, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	utils.SetSuccessJsonResponse(c, nil)
}
