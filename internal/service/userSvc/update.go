package userSvc

import (
	"context"
	"fmt"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/models"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
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
