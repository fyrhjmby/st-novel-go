package service

import (
	"context"
	"errors"
	"st-novel-go/src/ai/model"
	"st-novel-go/src/ai/provider"
	settingsDao "st-novel-go/src/settings/dao"
)

// Chat performs a non-streaming chat completion.
func Chat(ctx context.Context, apiKeyID uint, userID uint, messages []model.ChatMessage) (*model.ChatResponse, error) {
	// 1. Get API Key configuration and verify ownership
	apiKey, err := settingsDao.GetAPIKeyByID(apiKeyID, userID)
	if err != nil {
		return nil, errors.New("invalid API key ID or permission denied")
	}

	// 2. Get the provider instance from the factory
	aiProvider, err := provider.GetProvider(apiKey)
	if err != nil {
		return nil, err
	}

	// 3. Prepare the chat configuration
	chatConfig := model.ChatConfig{
		Model: apiKey.DefaultModel,
		// TODO: Allow user to override these in the request payload
		Temperature: 0.7,
		MaxTokens:   2048,
		Stream:      false,
	}

	// 4. Call the provider's Chat method
	return aiProvider.Chat(ctx, messages, chatConfig)
}

// StreamChat performs a streaming chat completion.
func StreamChat(ctx context.Context, apiKeyID uint, userID uint, messages []model.ChatMessage) (<-chan model.StreamResponse, error) {
	// 1. Get API Key configuration and verify ownership
	apiKey, err := settingsDao.GetAPIKeyByID(apiKeyID, userID)
	if err != nil {
		return nil, errors.New("invalid API key ID or permission denied")
	}

	// 2. Get the provider instance from the factory
	aiProvider, err := provider.GetProvider(apiKey)
	if err != nil {
		return nil, err
	}

	// 3. Prepare the chat configuration
	chatConfig := model.ChatConfig{
		Model:  apiKey.DefaultModel,
		Stream: true,
	}

	// 4. Call the provider's StreamChat method
	return aiProvider.StreamChat(ctx, messages, chatConfig)
}
