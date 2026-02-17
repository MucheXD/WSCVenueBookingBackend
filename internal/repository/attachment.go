package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type AttachmentEntity struct {
	ID int `gorm:"column:id;primaryKey"`
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
	var attachmentEntity AttachmentEntity
	attachmentEntity.BizType = bizType
	attachmentEntity.BizID = bizID
	attachmentEntity.BizIndex = bizIndex
	attachmentEntity.fromDomain(modelA)
	txDB := database.DB.Create(&attachmentEntity)
	if txDB.Error != nil {
		return txDB.Error
	}
	FileObjectLinked(modelA.FileToken)
	return nil
}

func (a *AttachmentEntity) fromDomain(modelA *models.Attachment) {
	a.BizFileType = modelA.BizFileType
	a.BizFileName = modelA.BizFileName
	a.FileToken = modelA.FileToken
}

func (a *AttachmentEntity) toDomain() models.Attachment {
	return models.Attachment{
		BizFileType: a.BizFileType,
		BizFileName: a.BizFileName,
		FileToken:   a.FileToken,
	}
}
