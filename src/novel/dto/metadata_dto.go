// Package dto defines data transfer objects used for communication between layers.
package dto

// UpdateMetadataPayload defines the structure for updating a novel's metadata.
type UpdateMetadataPayload struct {
	Title             *string        `json:"title"`
	Description       *string        `json:"description"`
	Cover             *string        `json:"cover"`
	Tags              *[]NovelTagDTO `json:"tags"`
	Status            *string        `json:"status"`
	ReferenceNovelIDs *[]string      `json:"referenceNovelIds"`
}
