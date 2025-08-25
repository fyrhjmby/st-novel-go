// st-novel-go/src/novel/dto/recent_dto.go
package dto

type RecentActivityItemDTO struct {
	ID             string `json:"id"`
	NovelID        string `json:"novelId"`
	NovelTitle     string `json:"novelTitle"`
	NovelCover     string `json:"novelCover"`
	EditedItemType string `json:"editedItemType"`
	EditedItemName string `json:"editedItemName"`
	EditedAt       string `json:"editedAt"`
}

type LogRecentAccessPayload struct {
	NovelID string `json:"novelId" binding:"required"`
}
