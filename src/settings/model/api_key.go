package model

import "gorm.io/gorm"

type ProviderType string

const (
	OpenAI ProviderType = "OpenAI"
	Gemini ProviderType = "Gemini"
	Claude ProviderType = "Claude"
)

type KeyStatus string

const (
	Enabled  KeyStatus = "启用"
	Disabled KeyStatus = "暂停"
)

// APIKey stores the configuration for an AI provider key.
// In a real application, the APIKey field should be encrypted in the database.
type APIKey struct {
	gorm.Model
	UserID       uint         `gorm:"not null;index" json:"user_id"`
	Provider     ProviderType `gorm:"type:varchar(50);not null" json:"provider"`
	Name         string       `gorm:"type:varchar(255);not null" json:"name"`
	APIKey       string       `gorm:"type:varchar(255);not null" json:"-"` // Omitted from JSON responses for security
	BaseURL      string       `gorm:"type:varchar(255)" json:"base_url"`
	DefaultModel string       `gorm:"type:varchar(100);not null" json:"default_model"`
	Status       KeyStatus    `gorm:"type:varchar(20);default:'启用'" json:"status"`
	Calls        uint         `gorm:"default:0" json:"calls"`
}
