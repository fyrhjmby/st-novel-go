// st-novel-go/src/novel/dao/history_dao.go
package dao

import (
	"gorm.io/gorm"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/model"
)

func CreateHistoryVersion(version *model.HistoryVersion) error {
	return database.DB.Create(version).Error
}

func GetHistoryForDocument(documentID string, userID uint) ([]model.HistoryVersion, error) {
	var versions []model.HistoryVersion
	err := database.DB.
		Where("document_id = ? AND user_id = ?", documentID, userID).
		Order("created_at DESC").
		Find(&versions).Error
	return versions, err
}

func FindHistoryVersionByID(versionID string, userID uint) (*model.HistoryVersion, error) {
	var version model.HistoryVersion
	err := database.DB.
		Where("id = ? AND user_id = ?", versionID, userID).
		First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func DeleteHistoryForDocument(tx *gorm.DB, documentID string) error {
	return tx.Where("document_id = ?", documentID).Delete(&model.HistoryVersion{}).Error
}
