package provider

import (
	"context"
	"st-novel-go/src/ai/model"
)

// AIProvider defines the interface for interacting with different AI language models.
// Each supported AI service (like OpenAI, Gemini, Claude) will have an adapter that implements this interface.
type AIProvider interface {
	// Chat sends a standard, non-streaming request and gets a complete response.
	Chat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (*model.ChatResponse, error)

	// StreamChat sends a request and returns a channel to receive real-time, streamed response chunks.
	StreamChat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (<-chan model.StreamResponse, error)
}
