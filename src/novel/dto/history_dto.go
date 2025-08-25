// st-novel-go/src/novel/dto/history_dto.go
package dto

type HistoryVersionDTO struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}
