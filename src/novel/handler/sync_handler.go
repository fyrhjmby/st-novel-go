// sync_handler.go — 同步端点：返回自指定时间戳以来的变更
package handler

import (
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUpdatesHandler GET /api/novels/:novelId/updates?since=<unix_timestamp>
// 返回该小说自 since 时间戳以来的所有变更，供前端离线同步使用
func GetUpdatesHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	// 验证用户对该小说的所有权
	_, err := dao.FindNovelByID(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Novel not found or permission denied")
		return
	}

	sinceStr := c.DefaultQuery("since", "0")
	sinceUnix, err := strconv.ParseInt(sinceStr, 10, 64)
	if err != nil {
		sinceUnix = 0
	}
	since := time.Unix(sinceUnix, 0)

	// 收集变更：目录（卷+章节）、派生内容、笔记、元数据
	updates := make(map[string][]interface{})
	deletes := make(map[string][]string)

	// 卷变更
	volumes, err := dao.GetVolumesByNovelID(novelID)
	if err == nil {
		var updatedVolumes []interface{}
		for _, v := range volumes {
			if v.UpdatedAt.After(since) {
				updatedVolumes = append(updatedVolumes, v)
			}
		}
		if len(updatedVolumes) > 0 {
			updates["volumes"] = updatedVolumes
		}
	}

	// 章节变更
	chapters, err := dao.GetChaptersByNovelID(novelID)
	if err == nil {
		var updatedChapters []interface{}
		for _, ch := range chapters {
			if ch.UpdatedAt.After(since) {
				updatedChapters = append(updatedChapters, ch)
			}
		}
		if len(updatedChapters) > 0 {
			updates["chapters"] = updatedChapters
		}
	}

	// 派生内容变更
	derived, err := dao.GetDerivedContentForNovel(novelID)
	if err == nil {
		var updatedDerived []interface{}
		for _, d := range derived {
			if d.UpdatedAt.After(since) {
				updatedDerived = append(updatedDerived, d)
			}
		}
		if len(updatedDerived) > 0 {
			updates["derived_content"] = updatedDerived
		}
	}

	// 笔记变更
	notes, err := dao.GetNotesForNovel(novelID)
	if err == nil {
		var updatedNotes []interface{}
		for _, n := range notes {
			if n.UpdatedAt.After(since) {
				updatedNotes = append(updatedNotes, n)
			}
		}
		if len(updatedNotes) > 0 {
			updates["notes"] = updatedNotes
		}
	}

	// 元数据变更
	novel, err := dao.FindNovelByID(novelID, userClaims.UserID)
	if err == nil && novel.UpdatedAt.After(since) {
		metadataDTO := map[string]interface{}{
			"id":          novel.ID.String(),
			"title":       novel.Title,
			"description": novel.Description,
			"cover":       novel.Cover,
			"tags":        novel.Tags,
			"status":      novel.Status,
			"category":    novel.Category,
		}
		updates["metadata"] = []interface{}{metadataDTO}
	}

	newTimestamp := time.Now().Unix()

	utils.Success(c, gin.H{
		"updates":      updates,
		"deletes":      deletes,
		"newTimestamp": newTimestamp,
	})
}
