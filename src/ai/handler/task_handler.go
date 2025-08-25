// st-novel-go/src/ai/handler/task_handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"st-novel-go/src/ai/dto"
	"st-novel-go/src/ai/service"
	"st-novel-go/src/middleware"
	"st-novel-go/src/utils"
)

func GetAIProvidersHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	providers, err := service.GetAIProvidersForEditor(userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, providers)
}

func StreamAITaskHandler(c *gin.Context) {
	var payload dto.StreamAITaskPayload
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

	eventChan, err := service.StreamAITask(c.Request.Context(), payload, userClaims.UserID)
	if err != nil {
		// Before streaming starts, we can send a normal error
		utils.Fail(c, err.Error())
		return
	}

	// Manually stream the response to match frontend parser
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			// Client disconnected
			return false
		case event, ok := <-eventChan:
			if !ok {
				// Channel closed, streaming is done
				return false
			}
			jsonData, err := json.Marshal(event)
			if err != nil {
				// Try to send an error event
				errorEvent, _ := json.Marshal(dto.TaskStreamEvent{Event: "error", Error: "Failed to marshal stream data"})
				fmt.Fprintf(w, "data: %s\n\n", errorEvent)
				return false // Stop streaming on marshalling error
			}

			// Write data in SSE format
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			return true // continue streaming
		}
	})
}
