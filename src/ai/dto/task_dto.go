// st-novel-go/src/ai/dto/task_dto.go
package dto

type AIProviderConfigDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Model       string `json:"model"`
	Temperature int    `json:"temperature"`
	MaxTokens   int    `json:"maxTokens"`
	Description string `json:"description"`
}

type StreamAITaskPayload struct {
	Prompt          string              `json:"prompt" binding:"required"`
	Config          AIProviderConfigDTO `json:"config" binding:"required"`
	TaskType        string              `json:"taskType" binding:"required"`
	SourceItemTitle string              `json:"sourceItemTitle"`
}

type TaskStreamEvent struct {
	Event   string `json:"event"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}
