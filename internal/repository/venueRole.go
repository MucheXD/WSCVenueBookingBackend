package repository

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"gorm.io/gorm"
)

type VenueRoleEntity struct {
	VAGID           int    `gorm:"column:vagid;primaryKey;autoIncrement"`
	RoleName        string `gorm:"column:role_name"`
	RoleDescription string `gorm:"column:role_description"`
}

func (VenueRoleEntity) TableName() string {
	return "venue_roles"
}

func ListVenueRoles() ([]*models.VenueRole, error) {
	var entities []VenueRoleEntity
	if err := database.DB.
		Model(&VenueRoleEntity{}).
		Order("vagid ASC").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	result := make([]*models.VenueRole, 0, len(entities))
	for _, entity := range entities {
		result = append(result, entity.toDomain())
	}
	return result, nil
}

func GetVenueRoleByID(vagid int) (*models.VenueRole, error) {
	var entity VenueRoleEntity
	if err := database.DB.
		Model(&VenueRoleEntity{}).
		Where("vagid = ?", vagid).
		Take(&entity).Error; err != nil {
		return nil, err
	}
	return entity.toDomain(), nil
}

func CreateVenueRoleTx(tx *gorm.DB, role *models.VenueRole) (int, error) {
	entity := VenueRoleEntity{}
	entity.fromDomain(role)
	if err := tx.Create(&entity).Error; err != nil {
		return 0, err
	}
	return entity.VAGID, nil
}

func UpdateVenueRoleTx(tx *gorm.DB, role *models.VenueRole) error {
	entity := VenueRoleEntity{}
	entity.fromDomain(role)
	return tx.Model(&VenueRoleEntity{}).
		Where("vagid = ?", role.VAGID).
		Updates(&entity).Error
}

func (v *VenueRoleEntity) fromDomain(role *models.VenueRole) {
	v.VAGID = role.VAGID
	v.RoleName = role.Name
	v.RoleDescription = role.Description
}

func (v *VenueRoleEntity) toDomain() *models.VenueRole {
	return &models.VenueRole{
		VAGID:       v.VAGID,
		Name:        v.RoleName,
		Description: v.RoleDescription,
	}
}
