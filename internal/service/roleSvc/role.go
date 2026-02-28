package roleSvc

import (
	"context"
	"fmt"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
	"gorm.io/gorm"
)

type VenueAccessEntryInput struct {
	VenueID       int
	AllowReserve  bool
	AllowApproval bool
	AllowEdit     bool
	AllowManage   bool
}

type CreateRoleInput struct {
	Name        string
	Description string
	Accesses    []VenueAccessEntryInput
}

type UpdateRoleInput struct {
	Name        *string
	Description *string
	Accesses    *[]VenueAccessEntryInput
}

func ListRoles(ctx context.Context) ([]*models.VenueRole, error) {
	_ = ctx
	roles, err := repository.ListVenueRoles()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRoleQueryInDB, err)
	}
	return roles, nil
}

func CreateRole(ctx context.Context, input CreateRoleInput) (int, error) {
	if input.Name == "" {
		return 0, ErrRoleNameRequired
	}
	if err := validateAccessEntries(input.Accesses); err != nil {
		return 0, err
	}

	vagid := 0
	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		createdVAGID, err := repository.CreateVenueRoleTx(tx, &models.VenueRole{
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		vagid = createdVAGID

		rows, err := repository.CreateVenueAccessEntriesTx(tx, vagid, toVenueAccessEntities(input.Accesses))
		if err != nil {
			return err
		}
		if rows != int64(len(input.Accesses)) || rows == 0 {
			return ErrRoleAccessGroupNotMatched
		}
		return nil
	})
	if err != nil {
		if err == ErrRoleAccessGroupNotMatched {
			return 0, err
		}
		return 0, fmt.Errorf("%w: %w", ErrRoleCreateInDB, err)
	}

	// 刷新权限缓存
	if err := venuePermission.RefreshVenueAccessCache(repository.NewVenueAccessLoader()); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrRoleQueryInDB, err)
	}

	return vagid, nil
}

func UpdateRole(ctx context.Context, vagid int, input UpdateRoleInput) error {
	if vagid <= 0 {
		return ErrRoleNotFound
	}
	if input.Name != nil && *input.Name == "" {
		return ErrRoleNameRequired
	}
	if input.Accesses != nil {
		if err := validateAccessEntries(*input.Accesses); err != nil {
			return err
		}
	}

	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existingRole, err := repository.GetVenueRoleByID(vagid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrRoleNotFound
			}
			return err
		}

		utils.UpdateFieldPtr(&existingRole.Name, input.Name)
		utils.UpdateFieldPtr(&existingRole.Description, input.Description)

		if err := repository.UpdateVenueRoleTx(tx, existingRole); err != nil {
			return err
		}

		if input.Accesses != nil {
			if err := repository.DeleteVenueAccessGroupByIDTx(tx, vagid); err != nil {
				return err
			}
			rows, err := repository.CreateVenueAccessEntriesTx(tx, vagid, toVenueAccessEntities(*input.Accesses))
			if err != nil {
				return err
			}
			if rows != int64(len(*input.Accesses)) {
				return ErrRoleAccessGroupNotMatched
			}
		}

		return nil
	})
	if err != nil {
		if err == ErrRoleNotFound || err == ErrRoleAccessGroupNotMatched {
			return err
		}
		return fmt.Errorf("%w: %w", ErrRoleUpdateInDB, err)
	}

	// 刷新权限缓存
	if input.Accesses != nil {
		if err := venuePermission.RefreshVenueAccessCache(repository.NewVenueAccessLoader()); err != nil {
			return fmt.Errorf("%w: %w", ErrRoleQueryInDB, err)
		}
	}

	return nil
}

func validateAccessEntries(accesses []VenueAccessEntryInput) error {
	for _, access := range accesses {
		if access.VenueID <= 0 {
			return ErrRoleVenueIDInvalid
		}
		exists, err := repository.VenueExists(access.VenueID)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrRoleQueryInDB, err)
		}
		if !exists {
			return ErrRoleVenueIDInvalid
		}
	}
	return nil
}

func toVenueAccessEntities(accesses []VenueAccessEntryInput) []repository.VenueAccessEntity {
	if len(accesses) == 0 {
		return nil
	}
	entities := make([]repository.VenueAccessEntity, 0, len(accesses))
	for _, access := range accesses {
		entities = append(entities, repository.VenueAccessEntity{
			VenueID:       access.VenueID,
			AllowReserve:  access.AllowReserve,
			AllowApproval: access.AllowApproval,
			AllowEdit:     access.AllowEdit,
			AllowManage:   access.AllowManage,
		})
	}
	return entities
}
