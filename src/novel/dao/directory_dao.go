// ..\st-novel-go\src\novel\dao\directory_dao.go
package dao

import (
	"errors"
	"gorm.io/gorm"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/model"
)

func CreateVolume(volume *model.Volume) error {
	return database.DB.Create(volume).Error
}

func CreateChapter(chapter *model.Chapter) error {
	return database.DB.Create(chapter).Error
}

func FindVolumeByID(volumeID string) (*model.Volume, error) {
	var volume model.Volume
	if err := database.DB.Preload("Chapters").First(&volume, "id = ?", volumeID).Error; err != nil {
		return nil, err
	}
	return &volume, nil
}

func FindChapterByID(chapterID string) (*model.Chapter, error) {
	var chapter model.Chapter
	if err := database.DB.First(&chapter, "id = ?", chapterID).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

func GetVolumesByNovelID(novelID string) ([]model.Volume, error) {
	var volumes []model.Volume
	err := database.DB.Where("novel_id = ?", novelID).Order("`order` ASC").Find(&volumes).Error
	return volumes, err
}

func GetChaptersByNovelID(novelID string) ([]model.Chapter, error) {
	var chapters []model.Chapter
	err := database.DB.Where("novel_id = ?", novelID).Order("`order` ASC").Find(&chapters).Error
	return chapters, err
}

func UpdateVolume(volume *model.Volume) error {
	return database.DB.Save(volume).Error
}

func UpdateChapter(chapter *model.Chapter) error {
	return database.DB.Save(chapter).Error
}

func DeleteVolume(volumeID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Find all chapters of the volume
		var chapters []model.Chapter
		if err := tx.Where("volume_id = ?", volumeID).Find(&chapters).Error; err != nil {
			return err
		}

		// Delete related data for each chapter
		for _, chapter := range chapters {
			if err := DeleteHistoryForDocument(tx, chapter.ID.String()); err != nil {
				return err
			}
			if err := DeleteDerivedContentForSource(tx, chapter.ID.String()); err != nil {
				return err
			}
		}

		// Delete all chapters of the volume
		if err := tx.Where("volume_id = ?", volumeID).Delete(&model.Chapter{}).Error; err != nil {
			return err
		}

		// Delete related data for the volume itself
		if err := DeleteHistoryForDocument(tx, volumeID); err != nil {
			return err
		}
		if err := DeleteDerivedContentForSource(tx, volumeID); err != nil {
			return err
		}

		// Finally, delete the volume
		result := tx.Where("id = ?", volumeID).Delete(&model.Volume{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("volume not found")
		}

		return nil
	})
}

func DeleteChapter(chapterID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Delete related history
		if err := DeleteHistoryForDocument(tx, chapterID); err != nil {
			return err
		}
		// Delete related derived content
		if err := DeleteDerivedContentForSource(tx, chapterID); err != nil {
			return err
		}

		// Delete the chapter itself
		result := tx.Delete(&model.Chapter{}, "id = ?", chapterID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("chapter not found")
		}
		return nil
	})
}

func UpdateVolumeOrder(novelID string, orderedIDs []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range orderedIDs {
			if err := tx.Model(&model.Volume{}).Where("id = ? AND novel_id = ?", id, novelID).Update("order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func UpdateChapterOrder(volumeID string, orderedIDs []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range orderedIDs {
			if err := tx.Model(&model.Chapter{}).Where("id = ? AND volume_id = ?", id, volumeID).Update("order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
