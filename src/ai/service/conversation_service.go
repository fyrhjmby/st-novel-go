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

func UpdateConversationTitle(id string, userID uint, title string) (*ConversationDTO, error) {
	conv, err := dao.FindConversationByID(id, userID)
	if err != nil {
		return nil, err
	}
	conv.Title = title
	if err := dao.UpdateConversation(conv); err != nil {
		return nil, err
	}
	return mapConversationToDTO(*conv), nil
}

func SaveConversationMessages(id string, userID uint, messages []ChatMessageDTO) error {
	conv, err := dao.FindConversationByID(id, userID)
	if err != nil {
		return err
	}
	msgBytes, _ := json.Marshal(messages)
	conv.Messages = msgBytes
	return dao.UpdateConversation(conv)
}

func DeleteConversation(id string, userID uint) error {
	return dao.DeleteConversation(id, userID)
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
