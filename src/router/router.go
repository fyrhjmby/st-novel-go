// 文件路径: st-novel-go/src/router/router.go
package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	aiRouter "st-novel-go/src/ai/router"
	"st-novel-go/src/config"
	novelRouter "st-novel-go/src/novel/router"
	settingsRouter "st-novel-go/src/settings/router"
	userRouter "st-novel-go/src/user/router"
	"time"
)

func SetupRouter() *gin.Engine {
	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for your frontend URL in production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Main API group
	api := r.Group("/api")
	{
		// Register all module routes here
		userRouter.RegisterAuthRoutes(api.Group("/auth"))
		userRouter.RegisterUserRoutes(api.Group("/users"))
		settingsRouter.RegisterSettingsRoutes(api) // This now correctly registers routes like /api/api-keys
		aiRouter.RegisterAIRoutes(api)             // Registers routes under /api/ai
		novelRouter.RegisterNovelRoutes(api)       // Registers routes for novel dashboard, trash, etc.
	}

	return r
}
