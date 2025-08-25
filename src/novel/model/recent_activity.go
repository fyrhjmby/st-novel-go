// st-novel-go/src/novel/model/recent_activity.go
package model

import (
	"github.com/google/uuid"
)

type RecentActivity struct {
	BaseModel
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	NovelID        uuid.UUID `gorm:"type:char(36);not null;index" json:"novel_id"`
	EditedItemType string    `gorm:"type:varchar(50)" json:"edited_item_type"`
	EditedItemID   string    `gorm:"type:varchar(255)" json:"edited_item_id"`
	EditedItemName string    `gorm:"type:varchar(255)" json:"edited_item_name"`
}
