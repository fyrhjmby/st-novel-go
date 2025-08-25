// 文件路径: st-novel-go/src/novel/dao/novel_dao.go
package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/model"
)

func GetNovelsByUserID(userID uint) ([]model.Novel, error) {
	var novels []model.Novel
	err := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("updated_at desc").Find(&novels).Error
	return novels, err
}

type ChapterCountResult struct {
	NovelID uuid.UUID
	Count   int64
}

func GetChapterCountsForNovels(novelIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	var results []ChapterCountResult
	err := database.DB.Model(&model.Chapter{}).
		Select("novel_id, count(*) as count").
		Where("novel_id IN ?", novelIDs).
		Group("novel_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int64)
	for _, result := range results {
		counts[result.NovelID] = result.Count
	}
	return counts, nil
}

func CreateNovel(novel *model.Novel) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(novel).Error; err != nil {
			return err
		}

		defaultVolume := model.Volume{
			NovelID: novel.ID,
			Title:   "第一卷",
			Content: "<h1>第一卷</h1>",
			Order:   1,
		}
		if err := tx.Create(&defaultVolume).Error; err != nil {
			return err
		}

		return nil
	})
}

func FindNovelByID(novelID string, userID uint) (model.Novel, error) {
	var novel model.Novel
	err := database.DB.Where("id = ? AND user_id = ?", novelID, userID).First(&novel).Error
	return novel, err
}

func UpdateNovel(novel *model.Novel) error {
	return database.DB.Save(novel).Error
}

func UpdateNovelJSONField(novelID string, userID uint, fieldName string, data interface{}) error {
	// Verify user ownership first
	_, err := FindNovelByID(novelID, userID)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for field %s: %w", fieldName, err)
	}

	return database.DB.Model(&model.Novel{}).Where("id = ?", novelID).Update(fieldName, jsonData).Error
}

func SoftDeleteNovelByID(novelID string, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", novelID, userID).Delete(&model.Novel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("novel not found or permission denied")
	}
	return nil
}

func GetTrashedNovelsByUserID(userID uint) ([]model.Novel, error) {
	var novels []model.Novel
	err := database.DB.Unscoped().Where("user_id = ? AND deleted_at IS NOT NULL", userID).Order("deleted_at desc").Find(&novels).Error
	return novels, err
}

func RestoreNovelByID(novelID string, userID uint) error {
	result := database.DB.Unscoped().Model(&model.Novel{}).Where("id = ? AND user_id = ?", novelID, userID).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("trashed novel not found or permission denied")
	}
	return nil
}

func PermanentlyDeleteNovelByID(novelID string, userID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Find the novel and its children's IDs
		var novel model.Novel
		if err := tx.Unscoped().Preload("Volumes.Chapters").Where("id = ? AND user_id = ?", novelID, userID).First(&novel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("trashed novel not found or permission denied")
			}
			return err
		}

		// 2. Collect all document IDs for history cleanup
		docIDs := []string{novel.ID.String()}
		for _, volume := range novel.Volumes {
			docIDs = append(docIDs, volume.ID.String())
			for _, chapter := range volume.Chapters {
				docIDs = append(docIDs, chapter.ID.String())
			}
		}

		// 3. Delete related HistoryVersions
		if err := tx.Where("document_id IN ?", docIDs).Delete(&model.HistoryVersion{}).Error; err != nil {
			return err
		}

		// 4. Delete related RecentActivities
		if err := tx.Where("novel_id = ?", novel.ID).Delete(&model.RecentActivity{}).Error; err != nil {
			return err
		}

		// 5. Permanently delete the novel. GORM's `OnDelete:CASCADE` will handle
		// Volumes, Chapters, DerivedContents, and Notes.
		result := tx.Unscoped().Delete(&novel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			// This case should ideally not be hit due to the initial check, but as a safeguard:
			return errors.New("failed to delete the novel record")
		}

		return nil
	})
}
