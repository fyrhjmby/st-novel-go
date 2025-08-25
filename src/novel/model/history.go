// st-novel-go/src/novel/model/history.go
package model

type HistoryVersion struct {
	BaseModel
	DocumentID   string `gorm:"type:varchar(255);not null;index" json:"document_id"`
	DocumentType string `gorm:"type:varchar(50);not null;index" json:"document_type"`
	Label        string `gorm:"type:varchar(255)" json:"label"`
	Content      string `gorm:"type:longtext" json:"content"`
	UserID       uint   `gorm:"not null;index" json:"user_id"`
}
