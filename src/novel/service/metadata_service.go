// 文件路径: st-novel-go/src/novel/service/metadata_service.go
package service

import (
	"encoding/json"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
)

func GetNovelMetadata(novelID string, userID uint) (*dto.NovelMetadataDTO, error) {
	novel, err := dao.FindNovelByID(novelID, userID)
	if err != nil {
		return nil, err
	}

	var tags []dto.NovelTagDTO
	_ = json.Unmarshal(novel.Tags, &tags)

	var refIDs []string
	_ = json.Unmarshal(novel.ReferenceNovelIDs, &refIDs)

	metadata := &dto.NovelMetadataDTO{
		ID:                novel.ID.String(),
		Title:             novel.Title,
		Description:       novel.Description,
		Cover:             novel.Cover,
		Tags:              tags,
		Status:            novel.Status,
		ReferenceNovelIDs: refIDs,
	}
	return metadata, nil
}

func UpdateNovelMetadata(novelID string, userID uint, payload dto.UpdateMetadataPayload) (*dto.NovelMetadataDTO, error) {
	novel, err := dao.FindNovelByID(novelID, userID)
	if err != nil {
		return nil, err
	}

	if payload.Title != nil {
		novel.Title = *payload.Title
	}
	if payload.Description != nil {
		novel.Description = *payload.Description
	}
	if payload.Cover != nil {
		novel.Cover = *payload.Cover
	}
	if payload.Status != nil {
		novel.Status = *payload.Status
	}
	if payload.Tags != nil {
		tagsJSON, _ := json.Marshal(payload.Tags)
		novel.Tags = tagsJSON
	}
	if payload.ReferenceNovelIDs != nil {
		refsJSON, _ := json.Marshal(payload.ReferenceNovelIDs)
		novel.ReferenceNovelIDs = refsJSON
	}

	if err := dao.UpdateNovel(&novel); err != nil {
		return nil, err
	}

	return GetNovelMetadata(novelID, userID)
}
