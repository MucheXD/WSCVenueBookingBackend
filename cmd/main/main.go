package main

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/database"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/logger"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/server"
)

func main() {
	logger.InitLogger()
	database.InitDatabase()
	server.InitServer()
}
