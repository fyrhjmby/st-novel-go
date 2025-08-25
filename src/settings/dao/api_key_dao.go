package dao

import (
	"st-novel-go/src/database"
	"st-novel-go/src/settings/model"
)

func CreateAPIKey(apiKey *model.APIKey) error {
	return database.DB.Create(apiKey).Error
}

func GetAPIKeyByID(id uint, userID uint) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&apiKey).Error
	return &apiKey, err
}

func GetAPIKeysByUserID(userID uint) ([]model.APIKey, error) {
	var apiKeys []model.APIKey
	err := database.DB.Where("user_id = ?", userID).Find(&apiKeys).Error
	return apiKeys, err
}

func UpdateAPIKey(apiKey *model.APIKey) error {
	// Save will update all fields of the model
	return database.DB.Save(apiKey).Error
}

func DeleteAPIKey(id uint, userID uint) error {
	// Ensure the user owns the key before deleting
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.APIKey{}).Error
}
