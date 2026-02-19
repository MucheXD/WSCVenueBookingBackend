package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type UserEntity struct {
	UID          string    `gorm:"column:uid;primaryKey"`
	PasswordHash string    `gorm:"column:password_hash"`
	PasswordSalt string    `gorm:"column:password_salt"`
	RegisteredAt time.Time `gorm:"column:registered_at"`
	Username     string    `gorm:"column:username"`
	SchoolID     string    `gorm:"column:school_id"`
	PhoneNumber  string    `gorm:"column:phone_number"`
	RealName     string    `gorm:"column:real_name"`
	PermMap      uint64    `gorm:"column:perm_map"`
	PermVAGID    int       `gorm:"column:perm_vagid"`
}

func (UserEntity) TableName() string {
	return "users"
}

func CreateNewUser(modelU *models.User) error {
	var userEntity UserEntity
	userEntity.fromDomain(modelU)
	if err := database.DB.Create(&userEntity).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(userID string) (*models.User, error) {
	var userEntity UserEntity
	txDB := database.DB.
		Model(&UserEntity{}).
		Where(&UserEntity{UID: userID}).
		Take(&userEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return userEntity.toDomain(), nil
}

func GetUniqueUserByUsername(username string) (*models.User, error) {
	retUsr, err := FoundUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if len(retUsr) <= 0 {
		return nil, nil
	}
	if len(retUsr) > 1 {
		slog.Warn("Multiple users found with the same username", "username", username)
		return nil, fmt.Errorf("Multiple user found by GetUniqueUserByUsername method, username: %s", username)
	}
	return retUsr[0], nil
}

func FoundUserByUsername(username string) ([]*models.User, error) {
	var userEntities []UserEntity
	txDB := database.DB.
		Model(&UserEntity{}).
		Where(&UserEntity{Username: username}).
		Find(&userEntities)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	users := make([]*models.User, 0)
	for _, userEntity := range userEntities {
		users = append(users, userEntity.toDomain())
	}
	return users, nil
}

func IsUsernameExists(username string) (bool, error) {
	var count int64
	txDB := database.DB.
		Model(&UserEntity{}).
		Where(&UserEntity{Username: username}).
		Count(&count)
	if txDB.Error != nil {
		return false, txDB.Error
	}
	return count > 0, nil
}

func DeleteUserByID(userID string) error {
	txDB := database.DB.
		Model(&UserEntity{}).
		Where(&UserEntity{UID: userID}).
		Delete(&UserEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}
func EditUser(modelU *models.User) error {
	var userEntity UserEntity
	userEntity.fromDomain(modelU)
	if err := database.DB.Model(&UserEntity{}).
		Where(&UserEntity{UID: userEntity.UID}).
		Updates(&userEntity).Error; err != nil {
		return err
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (u *UserEntity) fromDomain(modelU *models.User) {
	u.UID = modelU.UID
	u.PasswordHash = modelU.PasswordHash
	u.PasswordSalt = modelU.PasswordSalt
	u.RegisteredAt = modelU.RegisterTime
	u.Username = modelU.Username
	u.SchoolID = modelU.SchoolID
	u.PhoneNumber=modelU.PhoneNumber
	u.RealName = modelU.RealName
	u.PermMap = modelU.PermMap
	u.PermVAGID = modelU.PermVAGID
}

func (u *UserEntity) toDomain() *models.User {
	return &models.User{
		UID:          u.UID,
		PasswordHash: u.PasswordHash,
		PasswordSalt: u.PasswordSalt,
		RegisterTime: u.RegisteredAt,
		Username:     u.Username,
		SchoolID:     u.SchoolID,
		PhoneNumber:  u.PhoneNumber,
		RealName:     u.RealName,
		PermMap:      u.PermMap,
		PermVAGID:    u.PermVAGID,
	}
}
