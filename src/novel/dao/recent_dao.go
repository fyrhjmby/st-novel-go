// st-novel-go/src/novel/dao/recent_dao.go
package dao

import (
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
	var existing model.RecentActivity
	// Find if an identical activity for the same item already exists
	err := database.DB.
		Where("user_id = ? AND novel_id = ? AND edited_item_id = ?", activity.UserID, activity.NovelID, activity.EditedItemID).
		First(&existing).Error

	if err == nil {
		// If it exists, update its timestamp
		return database.DB.Model(&existing).Update("updated_at", activity.UpdatedAt).Error
	}

	// If not found, create a new one
	return database.DB.Create(activity).Error
}
