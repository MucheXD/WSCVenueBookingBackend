package venueCtrl

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/venueSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateVenueForm 更新场地表单（所有字段可选）
type UpdateVenueForm struct {
	Name            string           `json:"name"`
	BuildingID      int              `json:"building_id"`
	TypeID          int              `json:"type_id"`
	Description     string           `json:"description"`
	CoverImageToken string           `json:"cover_image_token"`
	Equipments      *json.RawMessage `json:"equipments"`
	ImagesToken     []string         `json:"images_token"`
	Capacity        int              `json:"capacity"`
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
	if req.Equipments != nil {
		updates.EquipmentsRaw = *req.Equipments
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
		if errors.Is(err, venueSvc.ErrVenueEquipmentsInvalid) {
			apiException.AbortWithException(c, apiException.ParamError, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	// 处理图片附件更新
	if req.ImagesToken != nil {
		err = database.DB.Transaction(func(tx *gorm.DB) error {

			// 移除旧图片附件
			if err := repository.SoftDeleteBizAttachmentsTx(tx, repository.AttachmentBizTypeVenue, []int{venueID}); err != nil {
				return err
			}

			// 创建新图片附件条目
			attachments := make([]models.Attachment, 0, len(req.ImagesToken))
			for idx, fileToken := range req.ImagesToken {
				attachments = append(attachments, models.Attachment{
					Index:       idx,
					FileToken:   fileToken,
					BizFileType: "image",
				})
			}

			return repository.CreateAttachmentsTx(tx, repository.AttachmentBizTypeVenue, venueID, attachments)
		})
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
	}

	utils.SetSuccessJsonResponse(c, map[string]string{"status": "updated"})
}
