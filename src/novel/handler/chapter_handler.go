package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func GetChaptersForVolumeHandler(c *gin.Context) {
	volumeID := c.Param("volumeId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)
	chapters, err := service.GetChaptersForVolume(volumeID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, chapters)
}

func GetChapterHandler(c *gin.Context) {
	chapterID := c.Param("chapterId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)
	chapter, err := service.GetChapter(chapterID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, chapter)
}

func CreateChapterHandler(c *gin.Context) {
	volumeID := c.Param("volumeId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)
	var payload dto.CreateChapterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}
	chapter, err := service.CreateChapter(volumeID, userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, chapter)
}

func UpdateChapterHandler(c *gin.Context) {
	chapterID := c.Param("chapterId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)
	var payload dto.UpdateChapterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}
	chapter, err := service.UpdateChapter(chapterID, userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, chapter)
}

func DeleteChapterHandler(c *gin.Context) {
	chapterID := c.Param("chapterId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)
	if err := service.DeleteChapter(chapterID, userClaims.UserID); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Chapter deleted successfully")
}

func UpdateChapterOrderHandler(c *gin.Context) {
	volumeID := c.Param("volumeId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)
	var payload dto.OrderPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}
	if err := service.UpdateChapterOrder(volumeID, userClaims.UserID, payload.OrderedChapterIDs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Chapter order updated successfully")
}
