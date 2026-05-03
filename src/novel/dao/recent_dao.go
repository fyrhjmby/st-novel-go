// st-novel-go/src/novel/dao/recent_dao.go
package dao

import (
	"gorm.io/gorm/clause"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/model"
)

const RecentItemsLimit = 50

func GetRecentActivitiesByUserID(userID uint) ([]model.RecentActivity, error) {
	var activities []model.RecentActivity
	err := database.DB.
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(RecentItemsLimit).
		Find(&activities).Error
	return activities, err
}

func LogOrUpdateRecentActivity(activity *model.RecentActivity) error {
	// 使用 OnConflict 原子化 upsert，避免 TOCTOU 竞态条件。
	// 当 (user_id, novel_id, edited_item_id) 组合已存在时，只更新 updated_at。
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "novel_id"}, {Name: "edited_item_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(activity).Error
}
