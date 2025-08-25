package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"st-novel-go/src/ai/model"
	settingsModel "st-novel-go/src/settings/model"
	"strings"
	"time"
)

// Gemini-specific structures
type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role"` // "user" or "model"
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// GeminiAdapter is an adapter for the Google Gemini API.
type GeminiAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewGeminiAdapter creates a new adapter for Gemini.
func NewGeminiAdapter(config *settingsModel.APIKey) *GeminiAdapter {
	baseURL := "https://generativelanguage.googleapis.com/v1beta/models"
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	return &GeminiAdapter{
		apiKey:  config.APIKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (a *GeminiAdapter) Chat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (*model.ChatResponse, error) {
	return nil, errors.New("Gemini non-streaming Chat is not implemented")
}

func (a *GeminiAdapter) StreamChat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (<-chan model.StreamResponse, error) {
	geminiContents := a.prepareMessages(messages)
	reqBody := geminiRequest{
		Contents: geminiContents,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/%s:streamGenerateContent?key=%s", a.baseURL, config.Model, a.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	outChan := make(chan model.StreamResponse)
	go a.processStream(resp, outChan)

	return outChan, nil
}

func (a *GeminiAdapter) prepareMessages(messages []model.ChatMessage) []geminiContent {
	var contents []geminiContent
	for _, msg := range messages {
		if msg.Role == "system" {
			// Gemini handles system prompts differently, often by prepending to the first user message.
			// For simplicity here, we'll treat it as a user message.
			// A more advanced implementation could use the `system_instruction` field.
			contents = append(contents, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: msg.Content}},
			})
		} else {
			role := "user"
			if msg.Role == "assistant" {
				role = "model"
			}
			contents = append(contents, geminiContent{
				Role:  role,
				Parts: []geminiPart{{Text: msg.Content}},
			})
		}
	}
	return contents
}

func (a *GeminiAdapter) processStream(resp *http.Response, outChan chan<- model.StreamResponse) {
	defer resp.Body.Close()
	defer close(outChan)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		var streamResp geminiResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Candidates) > 0 && len(streamResp.Candidates[0].Content.Parts) > 0 {
			outChan <- model.StreamResponse{
				Content: streamResp.Candidates[0].Content.Parts[0].Text,
				Done:    false,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		outChan <- model.StreamResponse{Error: "stream reading error: " + err.Error(), Done: true}
	} else {
		outChan <- model.StreamResponse{Done: true}
	}
}
