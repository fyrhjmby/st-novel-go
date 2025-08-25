package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"st-novel-go/src/ai/model"
	"st-novel-go/src/ai/service"
	"st-novel-go/src/middleware"
	"st-novel-go/src/utils"
)

type ChatPayload struct {
	APIKeyID uint                `json:"api_key_id" binding:"required"`
	Messages []model.ChatMessage `json:"messages" binding:"required"`
}

func GetConversationsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	conversations, err := service.GetConversations(userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to fetch conversations: "+err.Error())
		return
	}
	utils.Success(c, conversations)
}

func CreateConversationHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	conversation, err := service.CreateConversation(userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to create conversation: "+err.Error())
		return
	}
	utils.Success(c, conversation)
}

func StreamChatHandler(c *gin.Context) {
	var payload ChatPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	streamChan, err := service.StreamChat(c.Request.Context(), payload.APIKeyID, userClaims.UserID, payload.Messages)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case chunk, ok := <-streamChan:
			if !ok {
				return false
			}
			jsonData, err := json.Marshal(chunk)
			if err != nil {
				errorChunk := model.StreamResponse{Error: "Failed to marshal stream data"}
				errorData, _ := json.Marshal(errorChunk)
				fmt.Fprintf(w, "data: %s\n\n", errorData)
				return true
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			return true
		}
	})
}
