package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

const (
	AttachmentBizTypeVenue              = 1
	AttachmentBizTypeApplication        = 2
	AttachmentBizTypeApplicationComment = 3
)

type AttachmentEntity struct {
	ID        int            `gorm:"column:id;primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	// Business Context
	BizType  int `gorm:"column:biz_type"`
	BizID    int `gorm:"column:biz_id"`
	BizIndex int `gorm:"column:biz_index"`
	// Attachment Info
	BizFileType string `gorm:"column:biz_filetype"`
	BizFileName string `gorm:"column:biz_filename"`
	FileToken   string `gorm:"column:file_token"`
}

func (AttachmentEntity) TableName() string {
	return "attachments"
}

func GetAttachmentsByBiz(bizType, bizID int) ([]models.Attachment, error) {
	var attachmentEntities []AttachmentEntity
	txDB := database.DB.
		Model(&AttachmentEntity{}).
		Where(&AttachmentEntity{BizType: bizType, BizID: bizID}).
		Order("biz_index ASC").
		Find(&attachmentEntities)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	var modelAttachments []models.Attachment
	for _, entity := range attachmentEntities {
		modelAttachments = append(modelAttachments, entity.toDomain())
	}
	return modelAttachments, nil
}

func CreateAttachment(modelA *models.Attachment, bizType, bizID, bizIndex int) error {
	return CreateAttachmentWithTx(database.DB, modelA, bizType, bizID, bizIndex)
}

func CreateAttachmentWithTx(tx *gorm.DB, modelA *models.Attachment, bizType, bizID, bizIndex int) error {
	var attachmentEntity AttachmentEntity
	attachmentEntity.BizType = bizType
	attachmentEntity.BizID = bizID
	attachmentEntity.BizIndex = bizIndex
	attachmentEntity.fromDomain(modelA)
	txDB := tx.Create(&attachmentEntity)
	if txDB.Error != nil {
		return txDB.Error
	}
	if err := FileObjectLinked(modelA.FileToken); err != nil {
		return err
	}
	return nil
}

func SoftDeleteAttachmentsByBizWithTx(tx *gorm.DB, bizType, bizID int) error {
	var attachments []AttachmentEntity
	if err := tx.
		Model(&AttachmentEntity{}).
		Where("biz_type = ? AND biz_id = ?", bizType, bizID).
		Find(&attachments).Error; err != nil {
		return err
	}

	for _, attachment := range attachments {
		if err := FileObjectUnlinked(attachment.FileToken); err != nil {
			return err
		}
	}

	return tx.
		Model(&AttachmentEntity{}).
		Where("biz_type = ? AND biz_id = ?", bizType, bizID).
		Delete(&AttachmentEntity{}).Error
}

func (a *AttachmentEntity) fromDomain(modelA *models.Attachment) {
	a.BizIndex = modelA.Index
	a.BizFileType = modelA.BizFileType
	a.BizFileName = modelA.BizFileName
	a.FileToken = modelA.FileToken
}

func (a *AttachmentEntity) toDomain() models.Attachment {
	return models.Attachment{
		Index:       a.BizIndex,
		BizFileType: a.BizFileType,
		BizFileName: a.BizFileName,
		FileToken:   a.FileToken,
	}
}
