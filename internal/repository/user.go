package repository

import (
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
)

type UserEntity struct {
	UUID         int       `gorm:"column:uuid;primaryKey"`
	PasswordHash string    `gorm:"column:password_hash"`
	RegisteredAt time.Time `gorm:"column:registered_at"`
	Username     string    `gorm:"column:username"`
	SchoolID     string    `gorm:"column:school_id"`
	RealName     string    `gorm:"column:real_name"`
	PermType     string    `gorm:"column:perm_type"`
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
		Where(&UserEntity{UUID: userID}).
		Take(&userEntity)
	if txDB.Error != nil {
		return nil, txDB.Error
	}
	return userEntity.toDomain(), nil
}

func DeleteUserByID(userID int) error {
	txDB := database.DB.
		Where(&UserEntity{UUID: userID}).
		Delete(&UserEntity{})
	if txDB.Error != nil {
		return txDB.Error
	}
	return nil
}

// Entity-Domain Conversion (Private Methods)

func (u *UserEntity) fromDomain(modelU *models.User) {
	u.UUID = modelU.ID
	u.PasswordHash = modelU.PasswordHash
	u.RegisteredAt = modelU.CreateTime
	u.Username = modelU.Username
	u.SchoolID = modelU.SchoolID
	u.RealName = modelU.RealName
	u.PermType = modelU.PermType
	u.PermVAGID = modelU.PermVAGID
}

func (u *UserEntity) toDomain() *models.User {
	return &models.User{
		ID:           u.UUID,
		PasswordHash: u.PasswordHash,
		CreateTime:   u.RegisteredAt,
		Username:     u.Username,
		SchoolID:     u.SchoolID,
		RealName:     u.RealName,
		PermType:     u.PermType,
		PermVAGID:    u.PermVAGID,
	}
}
