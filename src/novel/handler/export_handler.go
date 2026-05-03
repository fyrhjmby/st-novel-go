// export_handler.go — 小说导出端点
package handler

import (
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/utils"

	"github.com/gin-gonic/gin"
)

// ExportNovelHandler POST /api/novels/:novelId/export
// 当前为 stub 实现，返回小说正文的纯文本拼接
func ExportNovelHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	novel, err := dao.FindNovelByID(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Novel not found or permission denied")
		return
	}

	volumes, err := dao.GetVolumesByNovelID(novelID)
	if err != nil {
		utils.Fail(c, "Failed to load volumes: "+err.Error())
		return
	}

	var exportContent string
	exportContent += "# " + novel.Title + "\n\n"

	for _, vol := range volumes {
		exportContent += "## " + vol.Title + "\n\n"
		chapters, _ := dao.GetChaptersByVolumeID(vol.ID.String())
		for _, ch := range chapters {
			exportContent += stripHTML(ch.Content) + "\n\n"
		}
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\""+novel.Title+".txt\"")
	c.String(200, exportContent)
}

func stripHTML(html string) string {
	// 简单去除 HTML 标签
	result := ""
	inTag := false
	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result += string(r)
		}
	}
	return result
}
