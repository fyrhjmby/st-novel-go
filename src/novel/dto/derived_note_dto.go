package dto

type CreateDerivedContentPayload struct {
	Type     string `json:"type" binding:"required"`
	SourceID string `json:"sourceId" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
}

type UpdateDerivedContentPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type CreateNotePayload struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type UpdateNotePayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
