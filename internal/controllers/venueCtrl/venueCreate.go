package venueCtrl

import (
	"errors"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
)

// CreateVenueForm 创建场地表单
type CreateVenueForm struct {
	Name            string   `json:"name" binding:"required"`
	BuildingID      int      `json:"building_id" binding:"required"`
	TypeID          int      `json:"type_id"`
	Description     string   `json:"description"`
	CoverImageToken string   `json:"cover_image_token"`
	ImagesToken     []string `json:"images_token"`
	Capacity        int      `json:"capacity"`
}

// CreateVenueHandler 创建新场地
// PUT /api/venue
func CreateVenueHandler(c *gin.Context) {
	var req CreateVenueForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	// 构造场地模型
	venue := &models.Venue{
		Name:            req.Name,
		BuildingID:      req.BuildingID,
		TypeID:          req.TypeID,
		Description:     req.Description,
		CoverImageToken: req.CoverImageToken,
		Capacity:        req.Capacity,
	}

	// 调用服务层创建场地
	venueID, err := venueSvc.CreateVenue(c.Request.Context(), venue)
	if err != nil {
		if errors.Is(err, venueSvc.ErrVenueNameRequired) ||
			errors.Is(err, venueSvc.ErrVenueBuildingRequired) ||
			errors.Is(err, venueSvc.ErrVenueBuildingInvalid) {
			apiException.AbortWithException(c, apiException.ParamError, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// TODO: 处理附件图片（images_token）
	// 这里可以调用 repository.CreateAttachment 来创建附件记录

	utils.SetSuccessJsonResponse(c, map[string]int{"venue_id": venueID})
}
