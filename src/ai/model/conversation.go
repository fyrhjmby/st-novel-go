package model

import (
	"gorm.io/datatypes"
	base_model "st-novel-go/src/novel/model"
)

type Conversation struct {
	base_model.BaseModel
	UserID   uint           `gorm:"not null;index" json:"user_id"`
	Title    string         `gorm:"type:varchar(255);not null" json:"title"`
	Summary  string         `gorm:"type:text" json:"summary"`
	Messages datatypes.JSON `gorm:"type:json" json:"messages"`
}
