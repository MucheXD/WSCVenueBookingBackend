package utils

import (
	"log"
	"time"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/repository"
)

func ScheduleRelease() {
	now := time.Now()
	utcNow:=now.UTC()
	utcNowStr:=utcNow.Format("2006-01-02 15:04:05")
	result := database.DB.Model(&repository.NotificationContentEntity{}).Where("status=? AND release_time<=?", 2, utcNowStr).Update("status", 1)
	err := result.Error
	if err != nil {
		log.Printf("Error upload notifications: %v\n", err)
		return
	}
}