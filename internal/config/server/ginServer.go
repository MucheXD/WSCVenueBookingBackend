package server

import (
	"log"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var GinEngine *gin.Engine

func InitServer() {
	GinEngine = gin.Default()
	c:= cron.New()
	_, err:= c.AddFunc("* * * * *", utils.ScheduleRelease)
	if err != nil {
		log.Fatal("Add tasks error:", err)
	}
	c.Start()
	defer c.Stop()
	
	initRouter()
	GinEngine.Run(":8080")
}
