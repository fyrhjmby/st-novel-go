// 文件: ..\st-novel-go\src\novel\dao\project_dao.go

// st-novel-go/src/novel/dao/project_dao.go
package dao

import (
	"gorm.io/gorm"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/model"
)

type FullProjectData struct {
	Novel           model.Novel
	DerivedContents []model.DerivedContent
	Notes           []model.Note
}

func GetFullNovelProject(novelID string, userID uint) (*FullProjectData, error) {
	var novel model.Novel
	if err := database.DB.
		Preload("Volumes", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC")
		}).
		Preload("Volumes.Chapters", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC")
		}).
		Where("id = ? AND user_id = ?", novelID, userID).
		First(&novel).Error; err != nil {
		return nil, err
	}

	var derivedContents []model.DerivedContent
	if err := database.DB.Where("novel_id = ?", novelID).Find(&derivedContents).Error; err != nil {
		return nil, err
	}

	var notes []model.Note
	if err := database.DB.Where("novel_id = ?", novelID).Order("updated_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}

	return &FullProjectData{
		Novel:           novel,
		DerivedContents: derivedContents,
		Notes:           notes,
	}, nil
}

func GetAllNovelProjectsForUser(userID uint) ([]model.Novel, error) {
	var novels []model.Novel
	if err := database.DB.
		Preload("Volumes.Chapters").
		Preload("DerivedContents").
		Preload("Notes").
		Where("user_id = ?", userID).
		Find(&novels).Error; err != nil {
		return nil, err
	}
	return novels, nil
}

func CreateNovelProjectWithData(novel *model.Novel) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the Novel record itself (which includes JSON fields)
		if err := tx.Create(novel).Error; err != nil {
			return err
		}

		// 2. Create Volume and Chapter records, linking them correctly
		for i := range novel.Volumes {
			novel.Volumes[i].NovelID = novel.ID
			if err := tx.Create(&novel.Volumes[i]).Error; err != nil {
				return err
			}
			for j := range novel.Volumes[i].Chapters {
				novel.Volumes[i].Chapters[j].NovelID = novel.ID
				novel.Volumes[i].Chapters[j].VolumeID = novel.Volumes[i].ID
				if err := tx.Create(&novel.Volumes[i].Chapters[j]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
