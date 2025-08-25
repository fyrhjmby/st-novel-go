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

// Claude-specific structures
type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model    string          `json:"model"`
	System   string          `json:"system,omitempty"`
	Messages []claudeMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type claudeStreamDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type claudeStreamContentBlockDelta struct {
	Type  string            `json:"type"` // content_block_delta
	Index int               `json:"index"`
	Delta claudeStreamDelta `json:"delta"`
}

// ClaudeAdapter is an adapter for the Anthropic Claude API.
type ClaudeAdapter struct {
	apiKey     string
	baseURL    string
	client     *http.Client
	apiVersion string
}

// NewClaudeAdapter creates a new adapter for Claude.
func NewClaudeAdapter(config *settingsModel.APIKey) *ClaudeAdapter {
	baseURL := "https://api.anthropic.com/v1"
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	return &ClaudeAdapter{
		apiKey:     config.APIKey,
		baseURL:    baseURL,
		client:     &http.Client{Timeout: 60 * time.Second},
		apiVersion: "2023-06-01",
	}
}

func (a *ClaudeAdapter) Chat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (*model.ChatResponse, error) {
	return nil, errors.New("Claude non-streaming Chat is not implemented")
}

func (a *ClaudeAdapter) StreamChat(ctx context.Context, messages []model.ChatMessage, config model.ChatConfig) (<-chan model.StreamResponse, error) {
	systemPrompt, claudeMsgs := a.prepareMessages(messages)
	reqBody := claudeRequest{
		Model:    config.Model,
		Messages: claudeMsgs,
		System:   systemPrompt,
		Stream:   true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", a.apiVersion)

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

func (a *ClaudeAdapter) prepareMessages(messages []model.ChatMessage) (string, []claudeMessage) {
	var systemPrompt string
	var claudeMsgs []claudeMessage
	for i, msg := range messages {
		if i == 0 && msg.Role == "system" {
			systemPrompt = msg.Content
			continue
		}
		claudeMsgs = append(claudeMsgs, claudeMessage{Role: msg.Role, Content: msg.Content})
	}
	return systemPrompt, claudeMsgs
}

func (a *ClaudeAdapter) processStream(resp *http.Response, outChan chan<- model.StreamResponse) {
	defer resp.Body.Close()
	defer close(outChan)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		var eventData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &eventData); err != nil {
			continue
		}

		eventType, ok := eventData["type"].(string)
		if !ok {
			continue
		}

		if eventType == "content_block_delta" {
			deltaBytes, _ := json.Marshal(eventData)
			var contentDelta claudeStreamContentBlockDelta
			if json.Unmarshal(deltaBytes, &contentDelta) == nil {
				if contentDelta.Delta.Type == "text_delta" {
					outChan <- model.StreamResponse{
						Content: contentDelta.Delta.Text,
						Done:    false,
					}
				}
			}
		} else if eventType == "message_stop" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		outChan <- model.StreamResponse{Error: "stream reading error: " + err.Error(), Done: true}
	} else {
		outChan <- model.StreamResponse{Done: true}
	}
}
