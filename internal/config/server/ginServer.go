package server

import "github.com/gin-gonic/gin"

var GinEngine *gin.Engine

func InitServer() {
	GinEngine = gin.Default()
	initRouter()
	GinEngine.Run(":8080")
}
