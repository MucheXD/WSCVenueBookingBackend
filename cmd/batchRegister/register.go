package batchRegister

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/service/userSvc"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
)

func BatchRegisterFromFile(filePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return err
	}
	defer file.Close()

	// 初始化数据库（如果未初始化）
	if database.DB == nil {
		fmt.Println("Error: Database not initialized")
		return fmt.Errorf("database not initialized")
	}

	// 使用 CSV 读取器读取文件
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1 // 允许可变数量的列（2 或 3）

	ctx := context.Background()
	successCount := 0
	failCount := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error reading CSV line: %v\n", err)
			return err
		}

		// 检查至少有 2 列
		if len(record) < 2 {
			fmt.Printf("Error: Invalid record format, expected at least 2 columns\n")
			failCount++
			continue
		}

		username := strings.TrimSpace(record[0])
		password := strings.TrimSpace(record[1])
		schoolID := "" // 第三列可选

		if password == "" {
			fmt.Printf("Error: Password cannot be empty for user '%s'\n", username)
			failCount++
			continue
		}

		// 如果有第三列，使用它作为学校 ID
		if len(record) >= 3 {
			schoolID = strings.TrimSpace(record[2])
		}

		// 生成 16 位随机密码盐
		passwordSalt := utils.GenerateRandomStringULN(16)

		// 计算密码哈希：SHA256(password + salt)
		passwordBytes := []byte(password + passwordSalt)
		hashBytes := sha256.Sum256(passwordBytes)
		passwordHash := hex.EncodeToString(hashBytes[:])

		// 构造用户模型
		user := &models.User{
			Username:     username,
			SchoolID:     schoolID,
			PasswordHash: passwordHash,
			PasswordSalt: passwordSalt,
		}

		// 调用服务层注册用户
		_, err = userSvc.RegisterUser(ctx, user)
		if err != nil {
			fmt.Printf("Failed to register user '%s': %v\n", username, err)
			failCount++
			continue
		}

		fmt.Printf("Successfully registered user: %s\n", username)
		successCount++
	}

	// 输出汇总信息
	fmt.Printf("\n========== Registration Summary ==========\n")
	fmt.Printf("Total Success: %d\n", successCount)
	fmt.Printf("Total Failed: %d\n", failCount)
	fmt.Printf("==========================================\n")

	return nil
}
