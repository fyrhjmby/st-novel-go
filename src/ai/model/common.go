package model

// ChatMessage represents a single message in a conversation.
type ChatMessage struct {
	Role    string `json:"role"` // "system", "user", or "assistant"
	Content string `json:"content"`
}

// ChatConfig holds configuration options for a chat request.
type ChatConfig struct {
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"`
}

// ChatResponse is the structure for a non-streaming response.
type ChatResponse struct {
	Content string `json:"content"`
}

// StreamResponse is the structure for a chunk in a streaming response.
// Event 字段用于前端 SSE 解析：前端根据 "chunk"/"done"/"error" 区分事件类型。
type StreamResponse struct {
	Event   string `json:"event,omitempty"`
	Content string `json:"content,omitempty"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}
