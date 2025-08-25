package router

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/user/handler"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
}

func RegisterUserRoutes(router *gin.RouterGroup) {
	// Apply auth middleware to all routes in this group
	router.Use(middleware.AuthMiddleware())

	// Basic user info, e.g., for header display
	router.GET("/me", handler.GetCurrentUser)

	// Comprehensive settings for the user settings page
	router.GET("/settings", handler.GetUserSettingsHandler)
	router.PUT("/settings", handler.UpdateUserSettingsHandler)
}
