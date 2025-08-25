// 文件: ai_handler.go

package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AIStreamEvent defines the structure for all streaming events.
type AIStreamEvent struct {
	EventType string      `json:"eventType"`
	Payload   interface{} `json:"payload"`
}

// sendStreamEvent is a helper to marshal and send a structured SSE event.
func sendStreamEvent(c *gin.Context, eventName string, payload interface{}) {
	eventData, err := json.Marshal(payload)
	if err != nil {
		// Log the error on the server, but don't crash the stream
		fmt.Printf("Error marshalling stream event payload: %v\n", err)
		return
	}
	c.SSEvent(eventName, string(eventData))
	c.Writer.Flush()
}

// registerAIRoutes 注册所有与AI模块相关的路由
func registerAIRoutes(rg *gin.RouterGroup) {
	chatGroup := rg.Group("/ai/chat")
	{
		chatGroup.GET("/conversations", getConversationsHandler)
		chatGroup.POST("/conversations", createConversationHandler)
		chatGroup.POST("/conversations/:id/messages", sendMessageHandler)
	}

	taskGroup := rg.Group("/ai/tasks")
	{
		taskGroup.POST("/stream", streamAITaskHandler)
	}
}

// streamAITaskHandler 处理AI任务的流式请求 - BEST PRACTICE IMPLEMENTATION
func streamAITaskHandler(c *gin.Context) {
	var req AIStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 设置SSE头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		// 1. Send a "start" event immediately
		sendStreamEvent(c, "start", gin.H{"message": "AI task started"})
		time.Sleep(200 * time.Millisecond) // Simulate initial processing

		// 2. Simulate the actual AI call and stream back chunks
		// In a real scenario, you would be iterating over a channel from your AI services
		responseMessage := fmt.Sprintf("为《%s》执行的“%s”任务已收到，正在处理... 这是模拟的流式响应。", req.SourceItemTitle, req.TaskType)
		words := strings.Fields(responseMessage)

		for _, word := range words {
			// In a real OpenAI call, you'd get tokens, not words. This simulates that.
			chunk := word + " "
			sendStreamEvent(c, "chunk", gin.H{"content": chunk})
			time.Sleep(100 * time.Millisecond)
		}

		// 3. (Optional) Simulate an error
		// if someCondition {
		//   sendStreamEvent(c, "error", gin.H{"error": "An unexpected error occurred during generation."})
		// 	 return false // Terminate the stream on error
		// }

		// 4. Send a "done" event to signal completion
		sendStreamEvent(c, "done", gin.H{"message": "Stream completed successfully"})

		// Return false to close the connection
		return false
	})
}

// getConversationsHandler and other handlers remain the same...

// getConversationsHandler 获取所有对话列表
func getConversationsHandler(c *gin.Context) {
	AIStore.chatConversationsMutex.RLock()
	defer AIStore.chatConversationsMutex.RUnlock()

	convs := make([]Conversation, 0, len(AIStore.ChatConversations))
	for _, conv := range AIStore.ChatConversations {
		convs = append(convs, conv)
	}
	sort.Slice(convs, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, convs[i].CreatedAt)
		t2, err2 := time.Parse(time.RFC3339, convs[j].CreatedAt)
		if err1 == nil && err2 == nil {
			return t1.After(t2)
		}
		return convs[i].CreatedAt > convs[j].CreatedAt
	})
	c.JSON(http.StatusOK, convs)
}

// createConversationHandler 创建一个新的空对话
func createConversationHandler(c *gin.Context) {
	AIStore.chatConversationsMutex.Lock()
	defer AIStore.chatConversationsMutex.Unlock()

	newID := "conv-" + uuid.New().String()
	newConv := Conversation{
		ID:        newID,
		Title:     "新对话",
		Summary:   "开始新的对话...",
		CreatedAt: time.Now().Format(time.RFC3339),
		Messages:  []ChatMessage{},
	}
	AIStore.ChatConversations[newID] = newConv
	c.JSON(http.StatusCreated, newConv)
}

// sendMessageHandler 在指定对话中发送消息并获取AI回复
func sendMessageHandler(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	AIStore.chatConversationsMutex.Lock()
	defer AIStore.chatConversationsMutex.Unlock()

	conv, ok := AIStore.ChatConversations[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	userMessage := ChatMessage{
		ID:        "msg-" + uuid.New().String(),
		Role:      "user",
		Content:   req.Content,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	conv.Messages = append(conv.Messages, userMessage)

	// 模拟AI回复
	aiResponse := ChatMessage{
		ID:        "msg-" + uuid.New().String(),
		Role:      "ai",
		Content:   "<p>这是对 “" + req.Content + "” 的模拟AI回复。</p><p>我可以帮助您进行头脑风暴、角色设计、情节构思或文本润色。请告诉我您需要什么帮助。</p>",
		Timestamp: time.Now().Add(1 * time.Second).Format(time.RFC3339),
	}
	conv.Messages = append(conv.Messages, aiResponse)

	// 如果是新对话，更新标题和摘要
	if len(conv.Messages) == 2 {
		if len(req.Content) > 20 {
			conv.Title = req.Content[:20] + "..."
		} else {
			conv.Title = req.Content
		}
		conv.Summary = conv.Title
	}

	AIStore.ChatConversations[id] = conv
	c.JSON(http.StatusOK, gin.H{"userMessage": userMessage, "aiResponse": aiResponse})
}
