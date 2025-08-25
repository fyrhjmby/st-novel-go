// 文件路径: st-novel-go/src/novel/dto/dashboard_dto.go
package dto

type NovelTagDTO struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

type NovelStatusDTO struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

type NovelDashboardItemDTO struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Cover       string         `json:"cover"`
	Status      NovelStatusDTO `json:"status"`
	Tags        []NovelTagDTO  `json:"tags"`
	Chapters    int            `json:"chapters"`
	LastUpdated string         `json:"lastUpdated"`
	Category    string         `json:"category"`
}

type CreateNovelPayload struct {
	Title    string `json:"title" binding:"required"`
	Synopsis string `json:"synopsis"`
	Category string `json:"category" binding:"required"`
}

type DeletedItemDTO struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Icon             string `json:"icon"`
	DeletedAt        string `json:"deletedAt"`
	RetentionDays    int    `json:"retentionDays"`
	RetentionPercent int    `json:"retentionPercent"`
}

type NovelCategoryDTO struct {
	Name string `json:"name"`
}
