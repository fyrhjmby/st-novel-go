package service

import (
	"errors"
	"fmt"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/model"
)

func CreateVersion(documentID, documentType, label, content string, userID uint) error {
	version := &model.HistoryVersion{
		DocumentID:   documentID,
		DocumentType: documentType,
		Label:        label,
		Content:      content,
		UserID:       userID,
	}
	return dao.CreateHistoryVersion(version)
}

func GetHistory(documentID string, userID uint) ([]dto.HistoryVersionDTO, error) {
	versions, err := dao.GetHistoryForDocument(documentID, userID)
	if err != nil {
		return nil, err
	}

	var dtoList []dto.HistoryVersionDTO
	for _, v := range versions {
		dtoList = append(dtoList, dto.HistoryVersionDTO{
			ID:        v.ID.String(),
			Label:     v.Label,
			Timestamp: formatTimeAgo(v.CreatedAt),
			Content:   v.Content,
		})
	}
	return dtoList, nil
}

func RestoreVersion(documentID, versionID string, userID uint) error {
	versionToRestore, err := dao.FindHistoryVersionByID(versionID, userID)
	if err != nil {
		return errors.New("history version not found or permission denied")
	}

	restoreLabel := fmt.Sprintf("从 %s 恢复", formatTimeAgo(versionToRestore.CreatedAt))

	switch versionToRestore.DocumentType {
	case "chapter":
		chapter, err := dao.FindChapterByID(documentID)
		if err != nil {
			return errors.New("target chapter not found")
		}
		if _, err := dao.FindNovelByID(chapter.NovelID.String(), userID); err != nil {
			return errors.New("permission denied for target chapter")
		}
		_ = CreateVersion(documentID, "chapter", "恢复前快照", chapter.Content, userID)
		chapter.Content = versionToRestore.Content
		if err := dao.UpdateChapter(chapter); err != nil {
			return err
		}
		return CreateVersion(documentID, "chapter", restoreLabel, chapter.Content, userID)

	case "volume":
		volume, err := dao.FindVolumeByID(documentID)
		if err != nil {
			return errors.New("target volume not found")
		}
		if _, err := dao.FindNovelByID(volume.NovelID.String(), userID); err != nil {
			return errors.New("permission denied for target volume")
		}
		_ = CreateVersion(documentID, "volume", "恢复前快照", volume.Content, userID)
		volume.Content = versionToRestore.Content
		if err := dao.UpdateVolume(volume); err != nil {
			return err
		}
		return CreateVersion(documentID, "volume", restoreLabel, volume.Content, userID)

	case "derived_content":
		item, err := dao.FindDerivedContentByID(documentID)
		if err != nil {
			return errors.New("target derived content not found")
		}
		if _, err := dao.FindNovelByID(item.NovelID.String(), userID); err != nil {
			return errors.New("permission denied for target derived content")
		}
		_ = CreateVersion(documentID, "derived_content", "恢复前快照", item.Content, userID)
		item.Content = versionToRestore.Content
		if err := dao.UpdateDerivedContent(item); err != nil {
			return err
		}
		return CreateVersion(documentID, "derived_content", restoreLabel, item.Content, userID)

	case "note":
		note, err := dao.FindNoteByID(documentID)
		if err != nil {
			return errors.New("target note not found")
		}
		if _, err := dao.FindNovelByID(note.NovelID.String(), userID); err != nil {
			return errors.New("permission denied for target note")
		}
		_ = CreateVersion(documentID, "note", "恢复前快照", note.Content, userID)
		note.Content = versionToRestore.Content
		if err := dao.UpdateNote(note); err != nil {
			return err
		}
		return CreateVersion(documentID, "note", restoreLabel, note.Content, userID)

	default:
		return fmt.Errorf("restoring document type '%s' is not supported", versionToRestore.DocumentType)
	}
}
