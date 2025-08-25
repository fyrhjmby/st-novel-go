// 文件路径: st-novel-go/src/novel/model/novel_models.go
package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Novel struct {
	BaseModel
	UserID             uint             `gorm:"not null;index" json:"user_id"`
	Title              string           `gorm:"type:varchar(255);not null" json:"title"`
	Description        string           `gorm:"type:text" json:"description"`
	Cover              string           `gorm:"type:text" json:"cover"`
	Tags               datatypes.JSON   `json:"tags"`
	Status             string           `gorm:"type:varchar(50);default:'编辑中'" json:"status"`
	Category           string           `gorm:"type:varchar(50)" json:"category"`
	ReferenceNovelIDs  datatypes.JSON   `json:"reference_novel_ids"`
	SettingsData       datatypes.JSON   `gorm:"type:json" json:"settings_data"`
	PlotCustomData     datatypes.JSON   `gorm:"type:json" json:"plot_custom_data"`
	AnalysisCustomData datatypes.JSON   `gorm:"type:json" json:"analysis_custom_data"`
	OthersCustomData   datatypes.JSON   `gorm:"type:json" json:"others_custom_data"`
	Volumes            []Volume         `gorm:"foreignKey:NovelID;constraint:OnDelete:CASCADE;" json:"volumes"`
	DerivedContents    []DerivedContent `gorm:"foreignKey:NovelID;constraint:OnDelete:CASCADE;" json:"derived_contents"`
	Notes              []Note           `gorm:"foreignKey:NovelID;constraint:OnDelete:CASCADE;" json:"notes"`
}

type Volume struct {
	BaseModel
	NovelID  uuid.UUID `gorm:"type:char(36);not null;index" json:"novel_id"`
	Title    string    `gorm:"type:varchar(255);not null" json:"title"`
	Content  string    `gorm:"type:longtext" json:"content"`
	Order    int       `gorm:"default:0" json:"order"`
	Chapters []Chapter `gorm:"foreignKey:VolumeID;constraint:OnDelete:CASCADE;" json:"chapters"`
}

type Chapter struct {
	BaseModel
	NovelID   uuid.UUID `gorm:"type:char(36);not null;index" json:"novel_id"`
	VolumeID  uuid.UUID `gorm:"type:char(36);not null;index" json:"volume_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	WordCount int       `gorm:"default:0" json:"word_count"`
	Content   string    `gorm:"type:longtext" json:"content"`
	Status    string    `gorm:"type:varchar(50);default:'editing'" json:"status"`
	Order     int       `gorm:"default:0" json:"order"`
}

type DerivedContent struct {
	BaseModel
	NovelID  uuid.UUID `gorm:"type:char(36);not null;index" json:"novel_id"`
	SourceID string    `gorm:"type:varchar(255);not null;index" json:"source_id"` // Can be VolumeID or ChapterID
	Type     string    `gorm:"type:varchar(50);not null" json:"type"`             // 'plot' or 'analysis'
	Title    string    `gorm:"type:varchar(255);not null" json:"title"`
	Content  string    `gorm:"type:longtext" json:"content"`
}

type Note struct {
	BaseModel
	NovelID uuid.UUID `gorm:"type:char(36);not null;index" json:"novel_id"`
	Title   string    `gorm:"type:varchar(255);not null" json:"title"`
	Content string    `gorm:"type:longtext" json:"content"`
}
