package venueCtrl

import (
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

// UpdateVenueForm 更新场地表单（所有字段可选）
type UpdateVenueForm struct {
	Name            string   `json:"name"`
	BuildingID      int      `json:"building_id"`
	TypeID          int      `json:"type_id"`
	Description     string   `json:"description"`
	CoverImageToken string   `json:"cover_image_token"`
	ImagesToken     []string `json:"images_token"`
	Capacity        int      `json:"capacity"`
}

// UpdateVenueHandler 更新场地信息
// PUT /api/venue/:venue_id
// 权限要求：对应场地的Edit权限 OR AllVenueEdit系统权限
func UpdateVenueHandler(c *gin.Context) {
	// 获取场地ID
	venueIDStr := c.Param("venue_id")
	venueID, err := strconv.Atoi(venueIDStr)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	// 权限检查：需要场地Edit权限或AllVenueEdit系统权限
	if !checkVenueEditPermission(c, venueID) {
		apiException.AbortWithException(c, apiException.VenuePermNotSatisfied)
		return
	}

	var req UpdateVenueForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	// 构造更新数据（只包含提供的字段）
	updates := &models.Venue{
		ID:              venueID,
		Name:            req.Name,
		BuildingID:      req.BuildingID,
		TypeID:          req.TypeID,
		Description:     req.Description,
		CoverImageToken: req.CoverImageToken,
		Capacity:        req.Capacity,
	}

	// 调用服务层更新场地
	err = venueSvc.UpdateVenue(c.Request.Context(), updates)
	if err != nil {
		if errors.Is(err, venueSvc.ErrVenueNotFound) {
			apiException.AbortWithException(c, apiException.NotFound, err)
			return
		}
		if errors.Is(err, venueSvc.ErrVenueBuildingInvalid) {
			apiException.AbortWithException(c, apiException.ParamError, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// TODO: 处理附件图片（images_token）的更新

	utils.SetSuccessJsonResponse(c, map[string]string{"status": "updated"})
}
