package dao

import (
	"st-novel-go/src/database"
	"st-novel-go/src/user/model"
)

func FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FindUserByID(id uint) (*model.User, error) {
	var user model.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(user *model.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func UpdateUser(user *model.User) error {
	// Save updates all fields of the model based on its primary key
	result := database.DB.Save(user)
	return result.Error
}
