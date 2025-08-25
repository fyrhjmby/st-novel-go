// 文件: ..\st-novel-go\src\novel\service\project_service.go
package service

import (
	"encoding/json"
	"math/rand"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/model"
)

func GetNovelProject(novelID string, userID uint) (*dto.NovelProjectDTO, error) {
	fullData, err := dao.GetFullNovelProject(novelID, userID)
	if err != nil {
		return nil, err
	}

	projectDTO := mapFullDataToProjectDTO(fullData)
	return &projectDTO, nil
}

func GetAllNovelProjects(userID uint) ([]dto.NovelProjectDTO, error) {
	novels, err := dao.GetAllNovelProjectsForUser(userID)
	if err != nil {
		return nil, err
	}

	var projectDTOs []dto.NovelProjectDTO
	for _, novel := range novels {
		fullData := &dao.FullProjectData{
			Novel:           novel,
			DerivedContents: novel.DerivedContents,
			Notes:           novel.Notes,
		}
		projectDTOs = append(projectDTOs, mapFullDataToProjectDTO(fullData))
	}

	return projectDTOs, nil
}

func CreateFullNovelProject(payload dto.CreateFullNovelProjectPayload, userID uint) (*dto.NovelProjectDTO, error) {
	novel, err := mapCreatePayloadToModel(payload, userID)
	if err != nil {
		return nil, err
	}

	if err := dao.CreateNovelProjectWithData(novel); err != nil {
		return nil, err
	}

	return GetNovelProject(novel.ID.String(), userID)
}

func ImportNovelProject(payload dto.CreateFullNovelProjectPayload, userID uint) (*dto.NovelProjectDTO, error) {
	novel, err := mapCreatePayloadToModel(payload, userID)
	if err != nil {
		return nil, err
	}

	if err := dao.CreateNovelProjectWithData(novel); err != nil {
		return nil, err
	}

	return GetNovelProject(novel.ID.String(), userID)
}

func mapCreatePayloadToModel(payload dto.CreateFullNovelProjectPayload, userID uint) (*model.Novel, error) {
	settingsJSON, err := json.Marshal(payload.SettingsData)
	if err != nil {
		return nil, err
	}
	plotJSON, err := json.Marshal(payload.PlotCustomData)
	if err != nil {
		return nil, err
	}
	analysisJSON, err := json.Marshal(payload.AnalysisCustomData)
	if err != nil {
		return nil, err
	}
	othersJSON, err := json.Marshal(payload.OthersCustomData)
	if err != nil {
		return nil, err
	}
	tagsJSON, _ := json.Marshal([]dto.NovelTagDTO{})

	novel := &model.Novel{
		UserID:             userID,
		Title:              payload.Title,
		Description:        payload.Description,
		Category:           payload.Category,
		Status:             "编辑中",
		Cover:              defaultCovers[rand.Intn(len(defaultCovers))],
		Tags:               tagsJSON,
		SettingsData:       settingsJSON,
		PlotCustomData:     plotJSON,
		AnalysisCustomData: analysisJSON,
		OthersCustomData:   othersJSON,
	}

	for _, volDTO := range payload.DirectoryData {
		volume := model.Volume{
			Title:   volDTO.Title,
			Content: volDTO.Content,
			Order:   volDTO.Order,
		}
		for _, chapDTO := range volDTO.Chapters {
			chapter := model.Chapter{
				Title:     chapDTO.Title,
				Content:   chapDTO.Content,
				WordCount: chapDTO.WordCount,
				Status:    chapDTO.Status,
				Order:     chapDTO.Order,
			}
			volume.Chapters = append(volume.Chapters, chapter)
		}
		novel.Volumes = append(novel.Volumes, volume)
	}

	return novel, nil
}

func mapFullDataToProjectDTO(data *dao.FullProjectData) dto.NovelProjectDTO {
	novel := data.Novel
	var tags []dto.NovelTagDTO
	_ = json.Unmarshal(novel.Tags, &tags)

	var refIDs []string
	_ = json.Unmarshal(novel.ReferenceNovelIDs, &refIDs)

	metadata := dto.NovelMetadataDTO{
		ID:                novel.ID.String(),
		Title:             novel.Title,
		Description:       novel.Description,
		Cover:             novel.Cover,
		Tags:              tags,
		Status:            novel.Status,
		ReferenceNovelIDs: refIDs,
	}

	var directoryData []dto.VolumeDTO
	for _, vol := range novel.Volumes {
		var chapters []dto.ChapterDTO
		for _, chap := range vol.Chapters {
			chapters = append(chapters, dto.ChapterDTO{
				ID:        chap.ID.String(),
				Type:      "chapter",
				VolumeID:  chap.VolumeID.String(),
				Title:     chap.Title,
				WordCount: chap.WordCount,
				Content:   chap.Content,
				Status:    chap.Status,
				Order:     chap.Order,
			})
		}
		directoryData = append(directoryData, dto.VolumeDTO{
			ID:       vol.ID.String(),
			Type:     "volume",
			Title:    vol.Title,
			Content:  vol.Content,
			Chapters: chapters,
			Order:    vol.Order,
		})
	}

	var derivedPlotData, derivedAnalysisData []dto.PlotAnalysisItemDTO
	for _, dc := range data.DerivedContents {
		item := dto.PlotAnalysisItemDTO{
			ID:       dc.ID.String(),
			Type:     dc.Type + "_item",
			SourceID: dc.SourceID,
			Title:    dc.Title,
			Content:  dc.Content,
		}
		if dc.Type == "plot" {
			derivedPlotData = append(derivedPlotData, item)
		} else if dc.Type == "analysis" {
			derivedAnalysisData = append(derivedAnalysisData, item)
		}
	}

	var noteData []dto.NoteItemDTO
	for _, note := range data.Notes {
		noteData = append(noteData, dto.NoteItemDTO{
			ID:        note.ID.String(),
			Type:      "note",
			Title:     note.Title,
			Content:   note.Content,
			Timestamp: note.UpdatedAt,
		})
	}

	unmarshalTreeNode := func(jsonData []byte) []dto.TreeNodeDTO {
		var data []dto.TreeNodeDTO
		if len(jsonData) > 0 {
			_ = json.Unmarshal(jsonData, &data)
		}
		if data == nil {
			return make([]dto.TreeNodeDTO, 0)
		}
		return data
	}

	unmarshalItemNode := func(jsonData []byte) []dto.ItemNodeDTO {
		var data []dto.ItemNodeDTO
		if len(jsonData) > 0 {
			_ = json.Unmarshal(jsonData, &data)
		}
		if data == nil {
			return make([]dto.ItemNodeDTO, 0)
		}
		return data
	}

	return dto.NovelProjectDTO{
		Metadata:            metadata,
		DirectoryData:       directoryData,
		SettingsData:        unmarshalTreeNode(novel.SettingsData),
		PlotCustomData:      unmarshalItemNode(novel.PlotCustomData),
		AnalysisCustomData:  unmarshalItemNode(novel.AnalysisCustomData),
		OthersCustomData:    unmarshalItemNode(novel.OthersCustomData),
		DerivedPlotData:     derivedPlotData,
		DerivedAnalysisData: derivedAnalysisData,
		NoteData:            noteData,
	}
}
