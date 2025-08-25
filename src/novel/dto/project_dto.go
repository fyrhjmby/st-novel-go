// 文件: ..\st-novel-go\src\novel\dto\project_dto.go

package dto

import "time"

type NovelMetadataDTO struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	Cover             string        `json:"cover"`
	Tags              []NovelTagDTO `json:"tags"`
	Status            string        `json:"status"`
	ReferenceNovelIDs []string      `json:"referenceNovelIds"`
}

type ChapterDTO struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	VolumeID  string `json:"volumeId"`
	Title     string `json:"title"`
	WordCount int    `json:"wordCount"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	Order     int    `json:"order"`
}

type VolumeDTO struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Content  string       `json:"content"`
	Chapters []ChapterDTO `json:"chapters"`
	Order    int          `json:"order"`
}

type TreeNodeDTO map[string]interface{}

type ItemNodeDTO struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Icon    string `json:"icon"`
	Content string `json:"content"`
}

type PlotAnalysisItemDTO struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	SourceID string `json:"sourceId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type NoteItemDTO struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type NovelProjectDTO struct {
	Metadata            NovelMetadataDTO      `json:"metadata"`
	DirectoryData       []VolumeDTO           `json:"directoryData"`
	SettingsData        []TreeNodeDTO         `json:"settingsData"`
	PlotCustomData      []ItemNodeDTO         `json:"plotCustomData"`
	AnalysisCustomData  []ItemNodeDTO         `json:"analysisCustomData"`
	DerivedPlotData     []PlotAnalysisItemDTO `json:"derivedPlotData"`
	DerivedAnalysisData []PlotAnalysisItemDTO `json:"derivedAnalysisData"`
	OthersCustomData    []ItemNodeDTO         `json:"othersCustomData"`
	NoteData            []NoteItemDTO         `json:"noteData"`
}

type CreateFullNovelProjectPayload struct {
	Title              string        `json:"title" binding:"required"`
	Description        string        `json:"description"`
	Category           string        `json:"category" binding:"required"`
	DirectoryData      []VolumeDTO   `json:"directoryData" binding:"required"`
	SettingsData       []TreeNodeDTO `json:"settingsData"`
	PlotCustomData     []ItemNodeDTO `json:"plotCustomData"`
	AnalysisCustomData []ItemNodeDTO `json:"analysisCustomData"`
	OthersCustomData   []ItemNodeDTO `json:"othersCustomData"`
}
