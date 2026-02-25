package utils

import (
	"log"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
)

func ScheduleRelease() {
	now := time.Now()
	result := database.DB.Model(&repository.NotificationEntity{}).Where("status=? AND release_time<=?", 2, now).Update("status", 1)
	err := result.Error
	if err != nil {
		log.Printf("Error upload notifications: %v\n", err)
		return
	}
}