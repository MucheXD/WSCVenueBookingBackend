package server

import "github.com/gin-gonic/gin"

func initRouter() {
	GinEngine.GET("/test", func(c *gin.Context) { c.String(200, "success") })
}
