package database

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	user := config.Config.GetString("database.user")
	pass := config.Config.GetString("database.password")
	port := config.Config.GetString("database.port")
	host := config.Config.GetString("database.host")
	DBname := config.Config.GetString("database.database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, DBname)
	slog.Debug("Connecting to database (gorm/mysql)", "dsn", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failed")
	}

	DB = db
}
