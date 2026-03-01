package userSvc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"gorm.io/gorm"
)

// 使用传入模型更新既有用户模型，若传入字段非零值则更新对应字段
// 注意：会更新所有可更新字段，对权限的判断请在 Controller 层进行
func UpdateUser(c context.Context, userID string, update models.User) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	utils.UpdateField(&user.Username, update.Username)
	utils.UpdateField(&user.RealName, update.RealName)
	utils.UpdateField(&user.PhoneNumber, update.PhoneNumber)
	utils.UpdateField(&user.SchoolID, update.SchoolID)
	utils.UpdateField(&user.PermVAGID, update.PermVAGID)

	err = repository.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrUpdateUserInDB, err)
	}
	return nil
}

// GetUserProfile returns user profile fields by user ID.
func GetUserProfile(c context.Context, userID string) (*models.User, error) {
	_ = c
	user, err := repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrQueryUserInDB, err)
	}
	return user, nil
}

func BatchUpdateUsersSystemPermission(c context.Context, userIDs []string, permMap uint64) error {
	_ = c

	cleanUserIDs := normalizeUserIDs(userIDs)
	if len(cleanUserIDs) == 0 {
		return nil
	}

	err := repository.BatchUpdateUsersPermMap(cleanUserIDs, permMap)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrUpdateUserInDB, err)
	}
	return nil
}

func BatchUpdateUsersVenueAccessGroup(c context.Context, userIDs []string, vagid int) error {
	_ = c

	if vagid < 0 {
		return ErrVenueAccessGroupInvalid
	}

	cleanUserIDs := normalizeUserIDs(userIDs)
	if len(cleanUserIDs) == 0 {
		return nil
	}

	if vagid > 0 {
		if _, err := repository.GetVenueRoleByID(vagid); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrVenueAccessGroupNotFound
			}
			return fmt.Errorf("%w:%w", ErrQueryUserInDB, err)
		}
	}

	err := repository.BatchUpdateUsersVAGID(cleanUserIDs, vagid)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrUpdateUserInDB, err)
	}
	return nil
}

func ListUsers(c context.Context, offset int, limit int) ([]*models.User, error) {
	_ = c

	users, err := repository.ListUsers(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrQueryUserInDB, err)
	}
	return users, nil
}

func normalizeUserIDs(userIDs []string) []string {
	if len(userIDs) == 0 {
		return nil
	}

	cleanUserIDs := make([]string, 0, len(userIDs))
	seen := make(map[string]struct{}, len(userIDs))
	for _, uid := range userIDs {
		trUID := strings.TrimSpace(uid)
		trUID = strings.ToUpper(trUID)
		if trUID == "" {
			continue
		}
		if _, exists := seen[trUID]; exists {
			continue
		}
		seen[trUID] = struct{}{}
		cleanUserIDs = append(cleanUserIDs, trUID)
	}

	return cleanUserIDs
}
