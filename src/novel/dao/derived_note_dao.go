// ..\st-novel-go\src\novel\dao\derived_note_dao.go
package dao

import (
	"errors"
	"gorm.io/gorm"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/model"
)

// DerivedContent DAO
func GetDerivedContentForNovel(novelID string) ([]model.DerivedContent, error) {
	var items []model.DerivedContent
	err := database.DB.Where("novel_id = ?", novelID).Find(&items).Error
	return items, err
}

func FindDerivedContentByID(itemID string) (*model.DerivedContent, error) {
	var item model.DerivedContent
	if err := database.DB.First(&item, "id = ?", itemID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateDerivedContent(item *model.DerivedContent) error {
	return database.DB.Create(item).Error
}

func UpdateDerivedContent(item *model.DerivedContent) error {
	return database.DB.Save(item).Error
}

func DeleteDerivedContent(itemID string) error {
	result := database.DB.Delete(&model.DerivedContent{}, "id = ?", itemID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("derived content not found")
	}
	return nil
}

func DeleteDerivedContentForSource(tx *gorm.DB, sourceID string) error {
	return tx.Where("source_id = ?", sourceID).Delete(&model.DerivedContent{}).Error
}

// Note DAO
func GetNotesForNovel(novelID string) ([]model.Note, error) {
	var notes []model.Note
	err := database.DB.Where("novel_id = ?", novelID).Order("updated_at DESC").Find(&notes).Error
	return notes, err
}

func FindNoteByID(noteID string) (*model.Note, error) {
	var note model.Note
	if err := database.DB.First(&note, "id = ?", noteID).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func CreateNote(note *model.Note) error {
	return database.DB.Create(note).Error
}

func UpdateNote(note *model.Note) error {
	return database.DB.Save(note).Error
}

func DeleteNote(noteID string) error {
	result := database.DB.Delete(&model.Note{}, "id = ?", noteID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("note not found")
	}
	return nil
}
