// ..\st-novel-go\src\novel\dto\directory_dto.go
package dto

type CreateVolumePayload struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

type UpdateVolumePayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type CreateChapterPayload struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Status  string `json:"status"`
	Order   int    `json:"order"`
}

type UpdateChapterPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *string `json:"status"`
}

type OrderPayload struct {
	OrderedVolumeIDs  []string `json:"orderedVolumeIds"`
	OrderedChapterIDs []string `json:"orderedChapterIds"`
}
