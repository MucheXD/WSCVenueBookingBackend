package main

import (
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/logger"
	"github.com/MucheXD/WSCVenueBookingBackend/internal/config/server"
)

func main() {
	logger.InitLogger()
	server.InitServer()
	// database.InitDatabase()
}
