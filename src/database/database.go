package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	aiModel "st-novel-go/src/ai/model"
	"st-novel-go/src/config"
	novelModel "st-novel-go/src/novel/model"
	settingsModel "st-novel-go/src/settings/model"
	userModel "st-novel-go/src/user/model"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	err = DB.AutoMigrate(
		&userModel.User{},
		&settingsModel.APIKey{},
		&aiModel.Conversation{},
		&novelModel.Novel{},
		&novelModel.Volume{},
		&novelModel.Chapter{},
		&novelModel.DerivedContent{},
		&novelModel.Note{},
		&novelModel.RecentActivity{},
		&novelModel.HistoryVersion{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}
	log.Println("Database schema migrated successfully.")

	seedAdminUser()
}

func seedAdminUser() {
	var user userModel.User
	err := DB.Where("email = ?", "admin@example.com").First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Admin user not found, creating one...")
			admin := userModel.User{
				Email:    "admin@example.com",
				Password: "123456",
				Name:     "Admin",
			}
			if err := DB.Create(&admin).Error; err != nil {
				log.Fatalf("Failed to create admin user: %v", err)
			}
			log.Println("Admin user created successfully.")
		} else {
			log.Fatalf("Failed to check for admin user: %v", err)
		}
	} else {
		log.Println("Admin user already exists.")
	}
}
