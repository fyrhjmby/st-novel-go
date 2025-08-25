package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/model"
	"st-novel-go/src/utils"
	"time"
)

var defaultCovers = []string{
	"/covers/default-1.jpg",
	"/covers/default-2.jpg",
	"/covers/default-3.jpg",
	"/covers/default-4.jpg",
	"/covers/default-5.jpg",
}

func GetNovels(userID uint) ([]dto.NovelDashboardItemDTO, error) {
	novels, err := dao.GetNovelsByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(novels) == 0 {
		return []dto.NovelDashboardItemDTO{}, nil
	}

	novelIDs := make([]uuid.UUID, len(novels))
	for i, n := range novels {
		novelIDs[i] = n.ID
	}

	chapterCounts, err := dao.GetChapterCountsForNovels(novelIDs)
	if err != nil {
		return nil, err
	}

	var dtoList []dto.NovelDashboardItemDTO
	for _, novel := range novels {
		count := chapterCounts[novel.ID]
		dtoList = append(dtoList, mapNovelToDashboardDTO(novel, int(count)))
	}
	return dtoList, nil
}

func CreateNovel(payload dto.CreateNovelPayload, claims *utils.Claims) (dto.NovelDashboardItemDTO, error) {
	tagsJSON, _ := json.Marshal([]dto.NovelTagDTO{})
	randomCover := defaultCovers[rand.Intn(len(defaultCovers))]

	novel := model.Novel{
		UserID:             claims.UserID,
		Title:              payload.Title,
		Description:        payload.Synopsis,
		Category:           payload.Category,
		Status:             "编辑中",
		Tags:               tagsJSON,
		Cover:              randomCover,
		SettingsData:       GetDefaultSettingsData(),
		PlotCustomData:     GetDefaultCustomData(),
		AnalysisCustomData: GetDefaultCustomData(),
		OthersCustomData:   GetDefaultCustomData(),
	}

	err := dao.CreateNovel(&novel)
	if err != nil {
		return dto.NovelDashboardItemDTO{}, err
	}

	// For a new novel, chapter count is 0
	return mapNovelToDashboardDTO(novel, 0), nil
}

func MoveToTrash(novelID string, userID uint) error {
	return dao.SoftDeleteNovelByID(novelID, userID)
}

func GetAvailableCategories() []string {
	return []string{
		"科幻",
		"奇幻",
		"悬疑",
		"恐怖",
		"都市",
		"言情",
		"历史",
	}
}

func mapNovelToDashboardDTO(novel model.Novel, chapterCount int) dto.NovelDashboardItemDTO {
	var tags []dto.NovelTagDTO
	if novel.Tags != nil {
		_ = json.Unmarshal(novel.Tags, &tags)
	}

	statusMap := map[string]dto.NovelStatusDTO{
		"编辑中": {Text: "编辑中", Class: "bg-blue-100 text-blue-800"},
		"已完结": {Text: "已完结", Class: "bg-green-100 text-green-800"},
	}
	statusDTO, ok := statusMap[novel.Status]
	if !ok {
		statusDTO = dto.NovelStatusDTO{Text: novel.Status, Class: "bg-gray-100 text-gray-800"}
	}

	return dto.NovelDashboardItemDTO{
		ID:          novel.ID.String(),
		Title:       novel.Title,
		Description: novel.Description,
		Cover:       novel.Cover,
		Status:      statusDTO,
		Tags:        tags,
		Chapters:    chapterCount,
		LastUpdated: formatTimeAgo(novel.UpdatedAt),
		Category:    novel.Category,
	}
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes())

	if days > 0 {
		return fmt.Sprintf("%d天前", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时前", hours)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d分钟前", minutes)
	}
	return "刚刚"
}
