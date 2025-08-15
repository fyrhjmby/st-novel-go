package main

// ChatMessage 定义了AI聊天中的单条消息结构
type ChatMessage struct {
	ID        string `json:"id"`
	Role      string `json:"role"` // "user" or "ai"
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// Conversation 定义了一个完整的对话，包含多条消息
type Conversation struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Summary   string        `json:"summary"`
	CreatedAt string        `json:"createdAt"`
	Messages  []ChatMessage `json:"messages"`
}

// AIProviderConfig 定义了AI提供商（模型）的配置
type AIProviderConfig struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
	Description string  `json:"description"`
}

// AIStreamRequest 定义了AI任务流式请求的结构
type AIStreamRequest struct {
	Prompt          string           `json:"prompt"`
	Config          AIProviderConfig `json:"config"`
	TaskType        string           `json:"taskType"`
	SourceItemTitle string           `json:"sourceItemTitle"`
}
