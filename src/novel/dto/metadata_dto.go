package dto

type UpdateMetadataPayload struct {
	Title             *string        `json:"title"`
	Description       *string        `json:"description"`
	Cover             *string        `json:"cover"`
	Tags              *[]NovelTagDTO `json:"tags"`
	Status            *string        `json:"status"`
	ReferenceNovelIDs *[]string      `json:"referenceNovelIds"`
}
