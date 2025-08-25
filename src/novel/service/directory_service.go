package service

import (
	"errors"
	"github.com/google/uuid"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/model"
)

func GetVolumes(novelID string, userID uint) ([]model.Volume, error) {
	if _, err := dao.FindNovelByID(novelID, userID); err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	return dao.GetVolumesByNovelID(novelID)
}

func GetChaptersForNovel(novelID string, userID uint) ([]model.Chapter, error) {
	if _, err := dao.FindNovelByID(novelID, userID); err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	return dao.GetChaptersByNovelID(novelID)
}

func GetChapter(chapterID string, userID uint) (*model.Chapter, error) {
	chapter, err := dao.FindChapterByID(chapterID)
	if err != nil {
		return nil, err
	}
	if _, err := dao.FindNovelByID(chapter.NovelID.String(), userID); err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	return chapter, nil
}

func CreateVolume(novelID string, userID uint, payload dto.CreateVolumePayload) (*model.Volume, error) {
	novel, err := dao.FindNovelByID(novelID, userID)
	if err != nil {
		return nil, errors.New("permission denied or novel not found")
	}
	volume := &model.Volume{
		NovelID: novel.ID,
		Title:   payload.Title,
		Content: payload.Content,
		Order:   payload.Order,
	}
	if err := dao.CreateVolume(volume); err != nil {
		return nil, err
	}
	return volume, nil
}

func CreateChapter(volumeID string, userID uint, payload dto.CreateChapterPayload) (*model.Chapter, error) {
	volume, err := dao.FindVolumeByID(volumeID)
	if err != nil {
		return nil, errors.New("volume not found")
	}
	if _, err := dao.FindNovelByID(volume.NovelID.String(), userID); err != nil {
		return nil, errors.New("permission denied")
	}
	chapter := &model.Chapter{
		NovelID:  volume.NovelID,
		VolumeID: volume.ID,
		Title:    payload.Title,
		Content:  payload.Content,
		Status:   payload.Status,
		Order:    payload.Order,
	}
	if err := dao.CreateChapter(chapter); err != nil {
		return nil, err
	}
	return chapter, nil
}

func UpdateVolume(volumeID string, userID uint, payload dto.UpdateVolumePayload) (*model.Volume, error) {
	volume, err := dao.FindVolumeByID(volumeID)
	if err != nil {
		return nil, errors.New("volume not found")
	}
	if _, err := dao.FindNovelByID(volume.NovelID.String(), userID); err != nil {
		return nil, errors.New("permission denied")
	}

	_ = CreateVersion(volume.ID.String(), "volume", "手动保存", volume.Content, userID)

	if payload.Title != nil {
		volume.Title = *payload.Title
	}
	if payload.Content != nil {
		volume.Content = *payload.Content
		volume.Title = SyncTitleFromContent(volume.Content, volume.Title)
	}
	if err := dao.UpdateVolume(volume); err != nil {
		return nil, err
	}

	_ = LogRecentEdit(userID, volume.NovelID, "outline", volume.ID.String(), volume.Title)

	return volume, nil
}

func UpdateChapter(chapterID string, userID uint, payload dto.UpdateChapterPayload) (*model.Chapter, error) {
	chapter, err := dao.FindChapterByID(chapterID)
	if err != nil {
		return nil, errors.New("chapter not found")
	}
	if _, err := dao.FindNovelByID(chapter.NovelID.String(), userID); err != nil {
		return nil, errors.New("permission denied")
	}

	_ = CreateVersion(chapter.ID.String(), "chapter", "手动保存", chapter.Content, userID)

	if payload.Title != nil {
		chapter.Title = *payload.Title
	}
	if payload.Content != nil {
		chapter.Content = *payload.Content
		chapter.Title = SyncTitleFromContent(chapter.Content, chapter.Title)
	}
	if payload.Status != nil {
		chapter.Status = *payload.Status
	}
	if err := dao.UpdateChapter(chapter); err != nil {
		return nil, err
	}

	_ = LogRecentEdit(userID, chapter.NovelID, "chapter", chapter.ID.String(), chapter.Title)

	return chapter, nil
}

func DeleteVolume(volumeID string, userID uint) error {
	volume, err := dao.FindVolumeByID(volumeID)
	if err != nil {
		return errors.New("volume not found")
	}
	if _, err := dao.FindNovelByID(volume.NovelID.String(), userID); err != nil {
		return errors.New("permission denied")
	}
	return dao.DeleteVolume(volume.ID.String())
}

func DeleteChapter(chapterID string, userID uint) error {
	chapter, err := dao.FindChapterByID(chapterID)
	if err != nil {
		return errors.New("chapter not found")
	}
	if _, err := dao.FindNovelByID(chapter.NovelID.String(), userID); err != nil {
		return errors.New("permission denied")
	}
	return dao.DeleteChapter(chapter.ID.String())
}

func UpdateVolumeOrder(novelID string, userID uint, orderedIDs []string) error {
	if _, err := dao.FindNovelByID(novelID, userID); err != nil {
		return errors.New("permission denied or novel not found")
	}
	return dao.UpdateVolumeOrder(novelID, orderedIDs)
}

func UpdateChapterOrder(volumeID string, userID uint, orderedIDs []string) error {
	volume, err := dao.FindVolumeByID(volumeID)
	if err != nil {
		return errors.New("volume not found")
	}
	if _, err := dao.FindNovelByID(volume.NovelID.String(), userID); err != nil {
		return errors.New("permission denied")
	}
	return dao.UpdateChapterOrder(volumeID, orderedIDs)
}

func MapUUIDs(uuids []uuid.UUID) []string {
	s := make([]string, len(uuids))
	for i, u := range uuids {
		s[i] = u.String()
	}
	return s
}
