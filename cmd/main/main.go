package main

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/logger"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/server"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils/venuePermission"
)

func main() {
	logger.InitLogger()
	database.InitDatabase()
	database.InitRedis()
	venuePermission.RefreshVenueAccessCache()
	server.InitServer()
}
