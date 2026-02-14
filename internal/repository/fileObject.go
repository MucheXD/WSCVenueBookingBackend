package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type FileObjectEntity struct {
	FID       int    `gorm:"column:fid;primaryKey"`
	FileToken string `gorm:"column:file_token"`
	FileHash  string `gorm:"column:file_hash"`
	FilePath  string `gorm:"column:file_path"`
	FileSize  int64  `gorm:"column:file_size"`
	LinkCount int    `gorm:"column:link_count"`
}

func (FileObjectEntity) TableName() string {
	return "file_objects"
}

// TODO: 文件落盘逻辑

// func createFileObject(modelF *models.FileObject) error {
// 	var entity FileObjectEntity
// 	entity.fromDomain(modelF)
// 	return database.DB.Create(&entity).Error
// }

// func getFileObjectByToken(fileToken string) (*models.FileObject, error) {
// 	var entity FileObjectEntity
// 	txDB := database.DB.
// 		Where(&FileObjectEntity{FileToken: fileToken}).
// 		Take(&entity)
// 	if txDB.Error != nil {
// 		return nil, txDB.Error
// 	}
// 	return entity.toDomain(), nil
// }

// func deleteFileObjectByToken(fileToken string) error {
// 	txDB := database.DB.
// 		Where(&FileObjectEntity{FileToken: fileToken}).
// 		Delete(&FileObjectEntity{})
// 	return txDB.Error
// }

func FileObjectLinked(fileToken string) error {
	return database.DB.
		Where(FileObjectEntity{FileToken: fileToken}).
		Update("link_count", gorm.Expr("link_count + ?", 1)).
		Error
}

func FileObjectUnlinked(fileToken string) error {
	return database.DB.
		Where(FileObjectEntity{FileToken: fileToken}).
		Update("link_count", gorm.Expr("GREATEST(link_count - 1, 0)")).
		Error
}

func (f *FileObjectEntity) fromDomain(modelF *models.FileObject) {
	f.FileToken = modelF.FileToken
	f.FileHash = modelF.FileHash
	f.FileSize = modelF.FileSize
	f.LinkCount = 0
}

func (f *FileObjectEntity) toDomain() *models.FileObject {
	return &models.FileObject{
		FileToken: f.FileToken,
		FileHash:  f.FileHash,
		FileSize:  f.FileSize,
	}
}
