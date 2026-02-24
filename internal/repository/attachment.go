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
	modelA.Index = bizIndex
	return CreateAttachmentsWithTx(tx, bizType, bizID, []models.Attachment{*modelA})
}

func CreateAttachmentsWithTx(tx *gorm.DB, bizType, bizID int, attachments []models.Attachment) error {
	if len(attachments) == 0 {
		return nil
	}

	// 批量转换并构建 entities 供附件表批量插入
	// 同时收集 fileTokens 供后续批量建引用调用
	entities := make([]AttachmentEntity, 0, len(attachments))
	fileTokens := make([]string, 0, len(attachments))
	for idx, attachment := range attachments {
		if attachment.Index < 0 {
			attachment.Index = idx
		}
		entity := AttachmentEntity{
			BizType:  bizType,
			BizID:    bizID,
			BizIndex: attachment.Index,
		}
		entity.fromDomain(&attachment)
		entities = append(entities, entity)
		fileTokens = append(fileTokens, attachment.FileToken)
	}

	if err := tx.Create(&entities).Error; err != nil {
		return err
	}

	return FileObjectLinkedBatchWithTx(tx, fileTokens)
}

func SoftDeleteAttachmentsByBizWithTx(tx *gorm.DB, bizType int, bizIDs []int) error {
	if len(bizIDs) == 0 {
		return nil
	}

	var attachments []AttachmentEntity
	if err := tx.
		Model(&AttachmentEntity{}).
		Where(&AttachmentEntity{BizType: bizType}).
		Where("biz_id IN ?", bizIDs).
		Select("file_token").
		Find(&attachments).Error; err != nil {
		return err
	}

	fileTokens := make([]string, 0, len(attachments))
	for _, attachment := range attachments {
		fileTokens = append(fileTokens, attachment.FileToken)
	}
	if err := FileObjectUnlinkedBatchWithTx(tx, fileTokens); err != nil {
		return err
	}

	return tx.
		Model(&AttachmentEntity{}).
		Where(&AttachmentEntity{BizType: bizType}).
		Where("biz_id IN ?", bizIDs).
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
