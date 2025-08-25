package service

import (
	"st-novel-go/src/settings/dao"
	"st-novel-go/src/settings/model"
	"strconv"
)

// GetModalProviders returns a simplified list of providers for the UI modal.
func GetModalProviders() []model.ModalProvider {
	return []model.ModalProvider{
		{Name: "OpenAI", ShortName: "GPT", Description: "行业领先模型"},
		{Name: "Claude", ShortName: "CLD", Description: "Anthropic出品"},
		{Name: "Gemini", ShortName: "GMN", Description: "Google强力支持"},
	}
}

// GetAPIProviders returns a detailed list of providers with user-specific stats.
func GetAPIProviders(userID uint) ([]model.ApiProvider, error) {
	// In a real application, this would involve more complex aggregation queries.
	// For now, we'll fetch all keys and process them in memory.
	allKeys, err := dao.GetAPIKeysByUserID(userID)
	if err != nil {
		return nil, err
	}

	providerStats := make(map[model.ProviderType]struct {
		ActiveKeys int
		TotalCalls uint
	})

	for _, key := range allKeys {
		stats := providerStats[key.Provider]
		stats.ActiveKeys++
		stats.TotalCalls += key.Calls
		providerStats[key.Provider] = stats
	}

	modalProviders := GetModalProviders()
	apiProviders := make([]model.ApiProvider, len(modalProviders))

	for i, p := range modalProviders {
		stats := providerStats[model.ProviderType(p.Name)]
		statusText := "未配置"
		if stats.ActiveKeys > 0 {
			statusText = strconv.Itoa(stats.ActiveKeys) + "个密钥"
		}

		apiProviders[i] = model.ApiProvider{
			Name:        p.Name,
			ShortName:   p.ShortName,
			Description: p.Description,
			StatusText:  statusText,
			ActiveKeys:  stats.ActiveKeys,
			TotalCalls:  formatCalls(stats.TotalCalls), // Formatting calls
		}
	}

	return apiProviders, nil
}

// formatCalls is a helper to format call counts for display.
func formatCalls(calls uint) string {
	if calls < 1000 {
		return strconv.Itoa(int(calls))
	}
	if calls < 10000 {
		return strconv.FormatFloat(float64(calls)/1000.0, 'f', 1, 64) + "k"
	}
	return strconv.Itoa(int(calls/1000)) + "k"
}
