package repository

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type UserEntity struct {
	UID          int       `gorm:"column:uid;primaryKey"`
	PasswordHash string    `gorm:"column:password_hash"`
	PasswordSalt string    `gorm:"column:password_salt"`
	RegisteredAt time.Time `gorm:"column:registered_at"`
	Username     string    `gorm:"column:username"`
	SchoolID     string    `gorm:"column:school_id"`
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

func GetUserByID(userID int) (*models.User, error) {
	var userEntity UserEntity
	txDB := database.DB.
		Where(&UserEntity{UID: userID}).
		Take(&userEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return userEntity.toDomain(), nil
}

func FoundUserByUsername(username string) ([]*models.User, error) {
	var userEntities []UserEntity
	txDB := database.DB.
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

func DeleteUserByID(userID int) error {
	txDB := database.DB.
		Where(&UserEntity{UID: userID}).
		Delete(&UserEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (u *UserEntity) fromDomain(modelU *models.User) {
	u.UID = modelU.UID
	u.PasswordHash = modelU.PasswordHash
	u.PasswordSalt = modelU.PasswordSalt
	u.RegisteredAt = modelU.CreateTime
	u.Username = modelU.Username
	u.SchoolID = modelU.SchoolID
	u.RealName = modelU.RealName
	u.PermMap = modelU.PermMap
	u.PermVAGID = modelU.PermVAGID
}

func (u *UserEntity) toDomain() *models.User {
	return &models.User{
		UID:          u.UID,
		PasswordHash: u.PasswordHash,
		PasswordSalt: u.PasswordSalt,
		CreateTime:   u.RegisteredAt,
		Username:     u.Username,
		SchoolID:     u.SchoolID,
		RealName:     u.RealName,
		PermMap:      u.PermMap,
		PermVAGID:    u.PermVAGID,
	}
}
