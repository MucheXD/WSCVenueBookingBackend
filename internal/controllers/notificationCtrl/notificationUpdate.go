package notificationCtrl

import (
	"errors"
	"strconv"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/notificationSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/apiException"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func UpdateNotificationHandler(c *gin.Context) {
	notificationIDStr := c.Param("notification_id")
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var req updateNotificationForm
	if err := c.ShouldBindJSON(&req); err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	updates, err := toUpdateNotificationModel(req)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	err = notificationSvc.UpdateNotification(c.Request.Context(), &updates,notificationID)
	if err != nil {
		if errors.Is(err, notificationSvc.ErrNotificationNotFound) {
			apiException.AbortWithException(c, apiException.NotFound, err)
			return
		}
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	if req.Attachments != nil {
		err = database.DB.Transaction(func(tx *gorm.DB) error {

			if err := repository.SoftDeleteBizAttachmentsTx(tx, repository.AttachmentBizTypeNotification, []int{notificationID}); err != nil {
				return err
			}

			attachments := make([]models.Attachment, 0, len(req.Attachments))
			for idx, attachment := range req.Attachments {
				attachments = append(attachments, models.Attachment{
					Index:       idx,
					FileToken:   attachment.FileToken,
					BizFileType: attachment.FileType,
					BizFileName: attachment.FileName,
				})
			}

			return repository.CreateAttachmentsTx(tx, repository.AttachmentBizTypeNotification, notificationID, attachments)
		})
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
	}

	utils.SetSuccessJsonResponse(c, map[string]string{"status": "updated"})
}