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

// OpenAI-specific request/response structures
type openAIRequest struct {
	Model       string              `json:"model"`
	Messages    []model.ChatMessage `json:"messages"`
	Stream      bool                `json:"stream"`
	Temperature float32             `json:"temperature,omitempty"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAIStreamChoice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
}

type openAIStreamResponse struct {
	Choices []openAIStreamChoice `json:"choices"`
}

// OpenAIAdapter is an adapter for the OpenAI API.
type OpenAIAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewOpenAIAdapter creates a new adapter for OpenAI.
func NewOpenAIAdapter(config *settingsModel.APIKey) *OpenAIAdapter {
	baseURL := "https://api.openai.com/v1"
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	return &OpenAIAdapter{
		apiKey:  config.APIKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (o *OpenAIAdapter) Chat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (*model.ChatResponse, error) {
	// Not implemented as streaming is the primary use case
	return nil, errors.New("OpenAI non-streaming Chat is not implemented")
}

func (o *OpenAIAdapter) StreamChat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (<-chan model.StreamResponse, error) {
	reqBody := openAIRequest{
		Model:    config.Model,
		Messages: messages,
		Stream:   true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	outChan := make(chan model.StreamResponse)
	go o.processStream(resp, outChan)

	return outChan, nil
}

func (o *OpenAIAdapter) processStream(resp *http.Response, outChan chan<- model.StreamResponse) {
	defer resp.Body.Close()
	defer close(outChan)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var streamResp openAIStreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			// handle error, maybe send an error chunk
			continue
		}

		if len(streamResp.Choices) > 0 {
			outChan <- model.StreamResponse{
				Content: streamResp.Choices[0].Delta.Content,
				Done:    false,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		// handle scanner error
		outChan <- model.StreamResponse{Error: "stream reading error: " + err.Error(), Done: true}
	} else {
		outChan <- model.StreamResponse{Done: true}
	}
}
