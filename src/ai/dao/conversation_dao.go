package dao

import (
	"st-novel-go/src/ai/model"
	"st-novel-go/src/database"
)

func CreateConversation(conv *model.Conversation) error {
	return database.DB.Create(conv).Error
}

func GetConversationsByUserID(userID uint) ([]model.Conversation, error) {
	var conversations []model.Conversation
	err := database.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&conversations).Error
	return conversations, err
}
