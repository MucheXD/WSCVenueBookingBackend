package repository

import (
	"strings"

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
	return FileObjectLinkedTx(database.DB, fileToken)
}

func FileObjectLinkedTx(tx *gorm.DB, fileToken string) error {
	if fileToken == "" {
		return nil
	}
	return tx.
		Model(&FileObjectEntity{}).
		Where(&FileObjectEntity{FileToken: fileToken}).
		Update("link_count", gorm.Expr("link_count + ?", 1)).
		Error
}

// 文件对象批量建引用
// 若 fileTokens 在列表重复出现，会累加计数
func FileObjectLinkedBatchTx(tx *gorm.DB, fileTokens []string) error {
	cases, args, ok := buildFileTokenCaseArgs(fileTokens)
	if !ok {
		return nil
	}

	query :=
		"UPDATE file_objects SET link_count = link_count + CASE file_token " +
			strings.Join(cases, " ") +
			" ELSE 0 END WHERE file_token IN ?"

	return tx.Exec(query, args...).Error
}

func FileObjectUnlinked(fileToken string) error {
	return FileObjectUnlinkedBatchTx(database.DB, []string{fileToken})
}

// 文件对象批量解引用
// 若 fileTokens 在列表重复出现，会累加计数
func FileObjectUnlinkedBatchTx(tx *gorm.DB, fileTokens []string) error {
	cases, args, ok := buildFileTokenCaseArgs(fileTokens)
	if !ok {
		return nil
	}

	query :=
		"UPDATE file_objects SET link_count = GREATEST(link_count - CASE file_token " +
			strings.Join(cases, " ") +
			" ELSE 0 END, 0) WHERE file_token IN ?"

	return tx.Exec(query, args...).Error
}

// buildFileTokenCaseArgs 将 fileTokens 结算为 CASE 分支与最终参数。
// 返回值中的 args 已包含末尾的 IN ? 参数（去重后的 tokens 列表）。
func buildFileTokenCaseArgs(fileTokens []string) ([]string, []any, bool) {
	if len(fileTokens) == 0 {
		return nil, nil, false
	}

	countMap := make(map[string]int)
	for _, token := range fileTokens {
		if token == "" {
			continue
		}
		countMap[token]++
	}
	if len(countMap) == 0 {
		return nil, nil, false
	}

	/*
		生成的 SQL 结构示例（假设 countMap 有两个条目 "A": 2, "B": 1）：
		最终阶段生成的 args = ["A", 2, "B", 1, ["A", "B"]]
		UPDATE file_objects
		SET link_count = GREATEST(link_count - CASE file_token
		    WHEN ? THEN ?  -- 对应 args[0], args[1] -> "A", 2
		    WHEN ? THEN ?  -- 对应 args[2], args[3] -> "B", 1
		    ELSE 0 END, 0)
		WHERE file_token IN ? -- 对应 args[4] -> ["A", "B"]
		* GREATEST 用于防止计数低于零
	*/

	cases := make([]string, 0, len(countMap))
	args := make([]any, 0, len(countMap)*2)
	tokens := make([]string, 0, len(countMap))
	for token, count := range countMap {
		cases = append(cases, "WHEN ? THEN ?")
		args = append(args, token, count)
		tokens = append(tokens, token)
	}
	args = append(args, tokens)
	return cases, args, true
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
