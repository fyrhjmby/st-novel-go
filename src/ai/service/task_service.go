// st-novel-go/src/ai/service/task_service.go
package service

import (
	"context"
	"fmt"
	"st-novel-go/src/ai/dto"
	"st-novel-go/src/ai/model"
	"st-novel-go/src/ai/provider"
	settingsDao "st-novel-go/src/settings/dao"
	settingsModel "st-novel-go/src/settings/model"
	"strconv"
)

func GetAIProvidersForEditor(userID uint) ([]dto.AIProviderConfigDTO, error) {
	apiKeys, err := settingsDao.GetAPIKeysByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get api keys: %w", err)
	}

	var providers []dto.AIProviderConfigDTO
	for _, key := range apiKeys {
		if key.Status == settingsModel.Enabled {
			providers = append(providers, dto.AIProviderConfigDTO{
				ID:          strconv.FormatUint(uint64(key.ID), 10),
				Name:        key.Name,
				Model:       key.DefaultModel,
				Description: fmt.Sprintf("Provider: %s, Model: %s", key.Provider, key.DefaultModel),
				// Default values, can be overridden by frontend
				Temperature: 70,
				MaxTokens:   2048,
			})
		}
	}
	return providers, nil
}

func StreamAITask(ctx context.Context, payload dto.StreamAITaskPayload, userID uint) (<-chan dto.TaskStreamEvent, error) {
	// Create a temporary APIKey model from the payload config to use the provider factory
	tempAPIKeyConfig := &settingsModel.APIKey{
		Provider:     settingsModel.ProviderType(payload.Config.Name),
		APIKey:       "", // The key is not needed here as it's assumed to be handled by the provider's base URL or managed service
		BaseURL:      "", // Assuming standard endpoints, or this would come from payload
		DefaultModel: payload.Config.Model,
	}

	// This is a simplified approach. A more robust system would fetch the actual API key from DB
	// based on the config ID and inject it. For now, we assume provider can be instantiated without a key (e.g. proxy)
	// Let's find the actual key for security
	apiKeyID, _ := strconv.ParseUint(payload.Config.ID, 10, 32)
	actualApiKey, err := settingsDao.GetAPIKeyByID(uint(apiKeyID), userID)
	if err != nil {
		return nil, fmt.Errorf("invalid or unauthorized api key id: %s", payload.Config.ID)
	}
	tempAPIKeyConfig.APIKey = actualApiKey.APIKey
	tempAPIKeyConfig.BaseURL = actualApiKey.BaseURL
	tempAPIKeyConfig.Provider = actualApiKey.Provider

	aiProvider, err := provider.GetProvider(tempAPIKeyConfig)
	if err != nil {
		return nil, err
	}

	chatConfig := model.ChatConfig{
		Model:       payload.Config.Model,
		Temperature: float32(payload.Config.Temperature) / 100.0,
		MaxTokens:   payload.Config.MaxTokens,
		Stream:      true,
	}

	messages := []model.ChatMessage{
		{Role: "user", Content: payload.Prompt},
	}

	providerChan, err := aiProvider.StreamChat(ctx, messages, chatConfig)
	if err != nil {
		return nil, err
	}

	// Create a new channel to transform the provider response to the task event format
	eventChan := make(chan dto.TaskStreamEvent)
	go func() {
		defer close(eventChan)
		for chunk := range providerChan {
			if chunk.Error != "" {
				eventChan <- dto.TaskStreamEvent{Event: "error", Error: chunk.Error}
				return // Stop on error
			}
			if chunk.Done {
				eventChan <- dto.TaskStreamEvent{Event: "done"}
				return
			}
			if chunk.Content != "" {
				eventChan <- dto.TaskStreamEvent{Event: "chunk", Content: chunk.Content}
			}
		}
	}()

	return eventChan, nil
}
