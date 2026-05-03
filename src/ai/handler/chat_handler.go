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
	APIKeyID       uint                `json:"api_key_id" binding:"required"`
	Messages       []model.ChatMessage `json:"messages" binding:"required"`
	ConversationID string              `json:"conversation_id"`
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

func UpdateConversationHandler(c *gin.Context) {
	convID := c.Param("id")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	conv, err := service.UpdateConversationTitle(convID, userClaims.UserID, payload.Title)
	if err != nil {
		utils.Fail(c, "Failed to update conversation: "+err.Error())
		return
	}
	utils.Success(c, conv)
}

func DeleteConversationHandler(c *gin.Context) {
	convID := c.Param("id")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	if err := service.DeleteConversation(convID, userClaims.UserID); err != nil {
		utils.Fail(c, "Failed to delete conversation: "+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Conversation deleted")
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

	var fullResponse string

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			// 流被中断时，保存已收到的部分消息
			if payload.ConversationID != "" && fullResponse != "" {
				saveChatMessages(payload, userClaims.UserID, fullResponse)
			}
			return false
		case chunk, ok := <-streamChan:
			if !ok {
				// 流正常结束，保存完整回复
				if payload.ConversationID != "" && fullResponse != "" {
					saveChatMessages(payload, userClaims.UserID, fullResponse)
				}
				return false
			}
			if chunk.Event == "chunk" {
				fullResponse += chunk.Content
			}
			jsonData, err := json.Marshal(chunk)
			if err != nil {
				errorChunk := model.StreamResponse{Event: "error", Error: "Failed to marshal stream data"}
				errorData, _ := json.Marshal(errorChunk)
				fmt.Fprintf(w, "data: %s\n\n", errorData)
				return true
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			return true
		}
	})
}

func saveChatMessages(payload ChatPayload, userID uint, aiResponse string) {
	msgs := make([]service.ChatMessageDTO, 0, len(payload.Messages)+1)
	for _, m := range payload.Messages {
		msgs = append(msgs, service.ChatMessageDTO{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	msgs = append(msgs, service.ChatMessageDTO{
		Role:    "ai",
		Content: aiResponse,
	})
	_ = service.SaveConversationMessages(payload.ConversationID, userID, msgs)
}
