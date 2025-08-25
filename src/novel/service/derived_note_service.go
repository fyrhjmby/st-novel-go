package service

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/model"
)

func getNovelIDFromSource(sourceID string) (uuid.UUID, error) {
	chapter, err := dao.FindChapterByID(sourceID)
	if err == nil {
		return chapter.NovelID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	volume, err := dao.FindVolumeByID(sourceID)
	if err == nil {
		return volume.NovelID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	return uuid.Nil, errors.New("source document not found")
}

func GetDerivedContentForNovel(novelID string, userID uint) ([]model.DerivedContent, error) {
	if _, err := dao.FindNovelByID(novelID, userID); err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	return dao.GetDerivedContentForNovel(novelID)
}

func CreateDerivedContent(userID uint, payload dto.CreateDerivedContentPayload) (*model.DerivedContent, error) {
	novelID, err := getNovelIDFromSource(payload.SourceID)
	if err != nil {
		return nil, err
	}

	if _, err := dao.FindNovelByID(novelID.String(), userID); err != nil {
		return nil, errors.New("permission denied for the associated novel")
	}

	item := &model.DerivedContent{
		NovelID:  novelID,
		Type:     payload.Type,
		SourceID: payload.SourceID,
		Title:    payload.Title,
		Content:  payload.Content,
	}
	if err := dao.CreateDerivedContent(item); err != nil {
		return nil, err
	}
	return item, nil
}

func UpdateDerivedContent(itemID string, userID uint, payload dto.UpdateDerivedContentPayload) (*model.DerivedContent, error) {
	item, err := dao.FindDerivedContentByID(itemID)
	if err != nil {
		return nil, errors.New("derived content not found")
	}
	if _, err := dao.FindNovelByID(item.NovelID.String(), userID); err != nil {
		return nil, errors.New("permission denied")
	}

	_ = CreateVersion(item.ID.String(), "derived_content", "手动保存", item.Content, userID)

	if payload.Title != nil {
		item.Title = *payload.Title
	}
	if payload.Content != nil {
		item.Content = *payload.Content
		item.Title = SyncTitleFromContent(item.Content, item.Title)
	}
	if err := dao.UpdateDerivedContent(item); err != nil {
		return nil, err
	}

	_ = LogRecentEdit(userID, item.NovelID, item.Type, item.ID.String(), item.Title)

	return item, nil
}

func DeleteDerivedContent(itemID string, userID uint) error {
	item, err := dao.FindDerivedContentByID(itemID)
	if err != nil {
		return errors.New("derived content not found")
	}
	if _, err := dao.FindNovelByID(item.NovelID.String(), userID); err != nil {
		return errors.New("permission denied")
	}
	return dao.DeleteDerivedContent(itemID)
}

func GetNotesForNovel(novelID string, userID uint) ([]model.Note, error) {
	if _, err := dao.FindNovelByID(novelID, userID); err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	return dao.GetNotesForNovel(novelID)
}

func CreateNote(novelID string, userID uint, payload dto.CreateNotePayload) (*model.Note, error) {
	novelUUID, err := uuid.Parse(novelID)
	if err != nil {
		return nil, errors.New("invalid novel ID")
	}
	if _, err := dao.FindNovelByID(novelUUID.String(), userID); err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	note := &model.Note{
		NovelID: novelUUID,
		Title:   payload.Title,
		Content: payload.Content,
	}
	if err := dao.CreateNote(note); err != nil {
		return nil, err
	}
	return note, nil
}

func UpdateNote(noteID string, userID uint, payload dto.UpdateNotePayload) (*model.Note, error) {
	note, err := dao.FindNoteByID(noteID)
	if err != nil {
		return nil, errors.New("note not found")
	}
	if _, err := dao.FindNovelByID(note.NovelID.String(), userID); err != nil {
		return nil, errors.New("permission denied")
	}

	_ = CreateVersion(note.ID.String(), "note", "手动保存", note.Content, userID)

	if payload.Title != nil {
		note.Title = *payload.Title
	}
	if payload.Content != nil {
		note.Content = *payload.Content
		note.Title = SyncTitleFromContent(note.Content, note.Title)
	}
	if err := dao.UpdateNote(note); err != nil {
		return nil, err
	}

	_ = LogRecentEdit(userID, note.NovelID, "note", note.ID.String(), note.Title)

	return note, nil
}

func DeleteNote(noteID string, userID uint) error {
	note, err := dao.FindNoteByID(noteID)
	if err != nil {
		return errors.New("note not found")
	}
	if _, err := dao.FindNovelByID(note.NovelID.String(), userID); err != nil {
		return errors.New("permission denied")
	}
	return dao.DeleteNote(noteID)
}
