package userSvc

import (
	"fmt"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
	"github.com/gin-gonic/gin"
)

func EditUser(c *gin.Context,phoneNumber string)error{
	userID:=c.GetString("UserID")
	user,err:=repository.GetUserByID(userID)
	if err != nil {
		return  err
	}

	user.PhoneNumber=phoneNumber
	err = repository.CreateNewUser(user)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrCreateUserInDB, err)
	}
	return nil


}