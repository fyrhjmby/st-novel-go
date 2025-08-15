// main.go
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ST-Novel Go Backend API
// @version 1.0
// @description This is a mock server for the ST-Novel frontend application, refactored into modules.

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	// 初始化所有模块的内存数据库
	InitDatabase()

	router := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // 调整为前端的端口
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// API路由组，所有业务API都在/api下
	apiGroup := router.Group("/api")
	{
		registerUserRoutes(apiGroup)
		registerNovelRoutes(apiGroup)
		registerAIRoutes(apiGroup)
		registerSettingRoutes(apiGroup)
	}

	// Swagger API文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	router.Run(":8080")
}
