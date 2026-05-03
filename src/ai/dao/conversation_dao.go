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

func FindConversationByID(id string, userID uint) (*model.Conversation, error) {
	var conv model.Conversation
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func UpdateConversation(conv *model.Conversation) error {
	return database.DB.Save(conv).Error
}

func DeleteConversation(id string, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Conversation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return database.DB.Error
	}
	return nil
}
