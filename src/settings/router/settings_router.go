package router

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/settings/handler"
)

// RegisterSettingsRoutes registers all routes related to the settings module.
func RegisterSettingsRoutes(router *gin.RouterGroup) {
	// All routes in this file are already under /api
	// No need to add /settings group here as it's handled in the main router

	// Provider routes - These are public or require auth depending on design
	// Let's assume they need auth to see user-specific stats
	providerGroup := router.Group("/api-providers")
	providerGroup.Use(middleware.AuthMiddleware())
	{
		providerGroup.GET("", handler.GetAPIProvidersHandler)
		providerGroup.GET("/modal", handler.GetModalProvidersHandler)
	}

	// API Key management routes
	apiKeysGroup := router.Group("/api-keys")
	apiKeysGroup.Use(middleware.AuthMiddleware())
	{
		apiKeysGroup.POST("", handler.CreateAPIKeyHandler)
		apiKeysGroup.GET("", handler.GetAPIKeysHandler)
		apiKeysGroup.PUT("/:id", handler.UpdateAPIKeyHandler)
		apiKeysGroup.DELETE("/:id", handler.DeleteAPIKeyHandler)
	}
}
