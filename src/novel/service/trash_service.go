// 文件路径: st-novel-go/src/novel/service/trash_service.go
package service

import (
	"fmt"
	"github.com/google/uuid"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"time"
)

const retentionPeriodDays = 30

func GetTrashedNovels(userID uint) ([]dto.DeletedItemDTO, error) {
	novels, err := dao.GetTrashedNovelsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var dtoList []dto.DeletedItemDTO
	for _, novel := range novels {
		now := time.Now()
		deletedAt := novel.DeletedAt.Time
		expiresAt := deletedAt.Add(retentionPeriodDays * 24 * time.Hour)
		remainingDuration := expiresAt.Sub(now)
		retentionDays := int(remainingDuration.Hours() / 24)
		if retentionDays < 0 {
			retentionDays = 0
		}

		totalDuration := expiresAt.Sub(deletedAt)
		retentionPercent := 0
		if totalDuration.Seconds() > 0 {
			retentionPercent = int((remainingDuration.Seconds() / totalDuration.Seconds()) * 100)
			if retentionPercent < 0 {
				retentionPercent = 0
			}
		}

		item := dto.DeletedItemDTO{
			ID:               novel.ID.String(),
			Name:             novel.Title,
			Type:             "小说",
			Icon:             `<svg fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><rect x="5" y="3" width="14" height="18" rx="2"/><path d="M9 7H15M9 11H15M9 15H13"/></svg>`,
			DeletedAt:        fmt.Sprintf("%s (%d天前)", deletedAt.Format("2006-01-02"), int(now.Sub(deletedAt).Hours()/24)),
			RetentionDays:    retentionDays,
			RetentionPercent: retentionPercent,
		}
		dtoList = append(dtoList, item)
	}

	return dtoList, nil
}

func RestoreNovel(novelID string, userID uint) (dto.NovelDashboardItemDTO, error) {
	err := dao.RestoreNovelByID(novelID, userID)
	if err != nil {
		return dto.NovelDashboardItemDTO{}, err
	}

	restoredNovel, err := dao.FindNovelByID(novelID, userID)
	if err != nil {
		return dto.NovelDashboardItemDTO{}, err
	}

	counts, err := dao.GetChapterCountsForNovels([]uuid.UUID{restoredNovel.ID})
	if err != nil {
		// Log the error but don't fail the request, default to 0 chapters
		counts = make(map[uuid.UUID]int64)
	}
	chapterCount := counts[restoredNovel.ID]

	return mapNovelToDashboardDTO(restoredNovel, int(chapterCount)), nil
}

func PermanentlyDeleteNovel(novelID string, userID uint) error {
	return dao.PermanentlyDeleteNovelByID(novelID, userID)
}
