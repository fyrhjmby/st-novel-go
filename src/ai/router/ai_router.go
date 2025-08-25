package router

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/ai/handler"
	"st-novel-go/src/middleware"
)

func RegisterAIRoutes(router *gin.RouterGroup) {
	aiGroup := router.Group("/ai")
	aiGroup.Use(middleware.AuthMiddleware())
	{
		chatGroup := aiGroup.Group("/chat")
		{
			chatGroup.GET("/conversations", handler.GetConversationsHandler)
			chatGroup.POST("/conversations", handler.CreateConversationHandler)
		}

		aiGroup.POST("/stream-chat", handler.StreamChatHandler)

		aiGroup.GET("/providers", handler.GetAIProvidersHandler)
		taskGroup := aiGroup.Group("/tasks")
		{
			taskGroup.POST("/stream", handler.StreamAITaskHandler)
		}
	}
}
