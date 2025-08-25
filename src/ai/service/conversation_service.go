package service

import (
	"encoding/json"
	"st-novel-go/src/ai/dao"
	"st-novel-go/src/ai/model"
	"time"
)

type ChatMessageDTO struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type ConversationDTO struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Summary   string           `json:"summary"`
	CreatedAt string           `json:"createdAt"`
	Messages  []ChatMessageDTO `json:"messages"`
}

func CreateConversation(userID uint) (*ConversationDTO, error) {
	emptyMessages, _ := json.Marshal([]ChatMessageDTO{})

	conv := &model.Conversation{
		UserID:   userID,
		Title:    "新的对话",
		Summary:  "开始一段新的对话...",
		Messages: emptyMessages,
	}

	if err := dao.CreateConversation(conv); err != nil {
		return nil, err
	}

	return mapConversationToDTO(*conv), nil
}

func GetConversations(userID uint) ([]ConversationDTO, error) {
	convs, err := dao.GetConversationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var dtoList []ConversationDTO
	for _, conv := range convs {
		dtoList = append(dtoList, *mapConversationToDTO(conv))
	}

	return dtoList, nil
}

func mapConversationToDTO(conv model.Conversation) *ConversationDTO {
	var messages []ChatMessageDTO
	_ = json.Unmarshal(conv.Messages, &messages)
	if messages == nil {
		messages = []ChatMessageDTO{}
	}

	return &ConversationDTO{
		ID:        conv.ID.String(),
		Title:     conv.Title,
		Summary:   conv.Summary,
		CreatedAt: conv.CreatedAt.Format(time.RFC3339),
		Messages:  messages,
	}
}
