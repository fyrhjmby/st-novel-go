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
type StreamResponse struct {
	Content string `json:"content,omitempty"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}
