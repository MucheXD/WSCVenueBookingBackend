package userCtrl

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

type UpdateUserSysPermDTO struct {
	UserID  []string `json:"uid" binding:"required"`
	SysPerm int      `json:"system_permission" binding:"required"`
}

func UpdateUserSysPermHandler(c *gin.Context) {
	var req UpdateUserSysPermDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	if len(req.UserID) == 0 || req.SysPerm < 0 {
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
