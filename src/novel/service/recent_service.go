package service

import (
	"errors"
	"github.com/google/uuid"
	"sort"
	"st-novel-go/src/database"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/model"
	"time"
)

func GetRecentItems(userID uint) ([]dto.RecentActivityItemDTO, error) {
	activities, err := dao.GetRecentActivitiesByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return []dto.RecentActivityItemDTO{}, nil
	}

	novelIDs := make(map[uuid.UUID]bool)
	for _, activity := range activities {
		novelIDs[activity.NovelID] = true
	}

	var ids []uuid.UUID
	for id := range novelIDs {
		ids = append(ids, id)
	}

	var novels []model.Novel
	if err := database.DB.Where("id IN ?", ids).Find(&novels).Error; err != nil {
		return nil, err
	}

	novelMap := make(map[uuid.UUID]model.Novel)
	for _, novel := range novels {
		novelMap[novel.ID] = novel
	}

	var dtoList []dto.RecentActivityItemDTO
	for _, activity := range activities {
		novel, ok := novelMap[activity.NovelID]
		if !ok {
			continue
		}

		dtoList = append(dtoList, dto.RecentActivityItemDTO{
			ID:             activity.ID.String(),
			NovelID:        activity.NovelID.String(),
			NovelTitle:     novel.Title,
			NovelCover:     novel.Cover,
			EditedItemType: activity.EditedItemType,
			EditedItemName: activity.EditedItemName,
			EditedAt:       activity.UpdatedAt.Format(time.RFC3339),
		})
	}

	return dtoList, nil
}

func LogRecentEdit(userID uint, novelID uuid.UUID, itemType string, itemID string, itemName string) error {
	novel, err := dao.FindNovelByID(novelID.String(), userID)
	if err != nil {
		return errors.New("novel not found or permission denied when logging recent activity")
	}

	activity := &model.RecentActivity{
		UserID:         userID,
		NovelID:        novel.ID,
		EditedItemType: itemType,
		EditedItemID:   itemID,
		EditedItemName: itemName,
		BaseModel:      model.BaseModel{UpdatedAt: time.Now()},
	}

	return dao.LogOrUpdateRecentActivity(activity)
}

func LogRecentAccess(payload dto.LogRecentAccessPayload, userID uint) (*dto.RecentActivityItemDTO, error) {
	novel, err := dao.FindNovelByID(payload.NovelID, userID)
	if err != nil {
		return nil, errors.New("novel not found or permission denied")
	}

	chapters, err := dao.GetChaptersByNovelID(payload.NovelID)
	if err != nil {
		return nil, errors.New("failed to retrieve chapters for novel")
	}

	var latestChapter model.Chapter
	if len(chapters) > 0 {
		sort.Slice(chapters, func(i, j int) bool {
			return chapters[i].UpdatedAt.After(chapters[j].UpdatedAt)
		})
		latestChapter = chapters[0]
	} else {
		volumes, _ := dao.GetVolumesByNovelID(payload.NovelID)
		if len(volumes) > 0 {
			activity := &model.RecentActivity{
				UserID:         userID,
				NovelID:        novel.ID,
				EditedItemType: "outline",
				EditedItemID:   volumes[0].ID.String(),
				EditedItemName: volumes[0].Title,
				BaseModel:      model.BaseModel{UpdatedAt: time.Now()},
			}
			if err := dao.LogOrUpdateRecentActivity(activity); err != nil {
				return nil, err
			}
			return &dto.RecentActivityItemDTO{
				ID:             activity.ID.String(),
				NovelID:        novel.ID.String(),
				NovelTitle:     novel.Title,
				NovelCover:     novel.Cover,
				EditedItemType: "outline",
				EditedItemName: volumes[0].Title,
				EditedAt:       activity.UpdatedAt.Format(time.RFC3339),
			}, nil
		}
		return nil, errors.New("cannot log access for an empty novel")
	}

	activity := &model.RecentActivity{
		UserID:         userID,
		NovelID:        novel.ID,
		EditedItemType: "chapter",
		EditedItemID:   latestChapter.ID.String(),
		EditedItemName: latestChapter.Title,
		BaseModel:      model.BaseModel{UpdatedAt: time.Now()},
	}

	if err := dao.LogOrUpdateRecentActivity(activity); err != nil {
		return nil, err
	}

	return &dto.RecentActivityItemDTO{
		ID:             activity.ID.String(),
		NovelID:        novel.ID.String(),
		NovelTitle:     novel.Title,
		NovelCover:     novel.Cover,
		EditedItemType: "chapter",
		EditedItemName: latestChapter.Title,
		EditedAt:       activity.UpdatedAt.Format(time.RFC3339),
	}, nil
}
