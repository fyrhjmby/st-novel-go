package service

import (
	"errors"
	"st-novel-go/src/settings/dao"
	"st-novel-go/src/settings/model"
	"strings"
)

type CreateAPIKeyPayload struct {
	Provider     model.ProviderType `json:"provider" binding:"required"`
	Name         string             `json:"name" binding:"required"`
	APIKey       string             `json:"apiKey" binding:"required"` // Matches frontend payload
	BaseURL      string             `json:"baseUrl"`
	DefaultModel string             `json:"model" binding:"required"`
}

type UpdateAPIKeyPayload struct {
	Name    *string `json:"name"`
	APIKey  *string `json:"key"` // Frontend sends 'key' for updates
	BaseURL *string `json:"baseUrl"`
	Model   *string `json:"model"`
}

// toAPIKeyResponse converts a database model to a frontend-friendly DTO.
func toAPIKeyResponse(apiKey model.APIKey) model.APIKeyResponse {
	keyFragment := ""
	if len(apiKey.APIKey) > 7 {
		keyFragment = apiKey.APIKey[:3] + "..." + apiKey.APIKey[len(apiKey.APIKey)-4:]
	}

	var providerShort string
	switch apiKey.Provider {
	case model.OpenAI:
		providerShort = "GPT"
	case model.Claude:
		providerShort = "CLD"
	case model.Gemini:
		providerShort = "GMN"
	}

	return model.APIKeyResponse{
		ID:            apiKey.ID,
		Provider:      string(apiKey.Provider),
		ProviderShort: providerShort,
		Name:          apiKey.Name,
		KeyFragment:   keyFragment,
		Model:         apiKey.DefaultModel,
		Calls:         formatCalls(apiKey.Calls),
		Status:        apiKey.Status,
		Created:       apiKey.CreatedAt.Format("2006-01-02"),
		BaseURL:       apiKey.BaseURL,
	}
}

func CreateAPIKey(payload CreateAPIKeyPayload, userID uint) (*model.APIKeyResponse, error) {
	apiKey := &model.APIKey{
		UserID:       userID,
		Provider:     payload.Provider,
		Name:         payload.Name,
		APIKey:       payload.APIKey,
		BaseURL:      payload.BaseURL,
		DefaultModel: payload.DefaultModel,
		Status:       model.Enabled, // Default status
	}

	if err := dao.CreateAPIKey(apiKey); err != nil {
		return nil, err
	}
	response := toAPIKeyResponse(*apiKey)
	return &response, nil
}

func GetAPIKeys(userID uint) ([]model.APIKeyResponse, error) {
	keys, err := dao.GetAPIKeysByUserID(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]model.APIKeyResponse, len(keys))
	for i, key := range keys {
		responses[i] = toAPIKeyResponse(key)
	}
	return responses, nil
}

func UpdateAPIKey(id uint, userID uint, payload UpdateAPIKeyPayload) (*model.APIKeyResponse, error) {
	apiKey, err := dao.GetAPIKeyByID(id, userID)
	if err != nil {
		return nil, errors.New("API key not found or you don't have permission")
	}

	if payload.Name != nil {
		apiKey.Name = *payload.Name
	}
	if payload.APIKey != nil && strings.TrimSpace(*payload.APIKey) != "" {
		apiKey.APIKey = *payload.APIKey
	}
	if payload.BaseURL != nil {
		apiKey.BaseURL = *payload.BaseURL
	}
	if payload.Model != nil {
		apiKey.DefaultModel = *payload.Model // <--- FIX: Changed from apiKey.Model to apiKey.DefaultModel
	}

	if err := dao.UpdateAPIKey(apiKey); err != nil {
		return nil, err
	}
	response := toAPIKeyResponse(*apiKey)
	return &response, nil
}

func DeleteAPIKey(id uint, userID uint) (*model.APIKeyResponse, error) {
	apiKey, err := dao.GetAPIKeyByID(id, userID)
	if err != nil {
		return nil, errors.New("API key not found or permission denied")
	}

	if err := dao.DeleteAPIKey(id, userID); err != nil {
		return nil, err
	}

	response := toAPIKeyResponse(*apiKey)
	return &response, nil
}
