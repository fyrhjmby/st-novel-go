package provider

import (
	"errors"
	settingsModel "st-novel-go/src/settings/model"
)

// GetProvider is a factory function that returns an instance of the requested AIProvider.
// It takes the APIKey model to configure the adapter.
func GetProvider(apiKey *settingsModel.APIKey) (AIProvider, error) {
	if apiKey == nil {
		return nil, errors.New("APIKey configuration cannot be nil")
	}

	switch apiKey.Provider {
	case settingsModel.OpenAI:
		return NewOpenAIAdapter(apiKey), nil
	case settingsModel.Gemini:
		return NewGeminiAdapter(apiKey), nil
	case settingsModel.Claude:
		return NewClaudeAdapter(apiKey), nil
	default:
		return nil, errors.New("unknown AI provider type")
	}
}
