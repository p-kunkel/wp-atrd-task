package mappings

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func CreateUrlMappings() *gin.Engine {
	server := gin.Default()

	basicPath := server.Group("v1")

	basicPath.POST("/secret", controllers.AddSecret)
	basicPath.GET("/secret/:hash", controllers.GetSecret)
	return server
}

func RunServer() error {
	return CreateUrlMappings().Run()
}
