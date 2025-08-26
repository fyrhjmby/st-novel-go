package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func MoveNovelToTrashHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	if err := service.MoveToTrash(novelID, userClaims.UserID); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Novel moved to trash")
}

func GetTrashedNovelsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	items, err := service.GetTrashedNovels(userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to fetch trashed items")
		return
	}
	utils.Success(c, items)
}

func RestoreNovelHandler(c *gin.Context) {
	itemID := c.Param("itemId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	restoredNovel, err := service.RestoreNovel(itemID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, restoredNovel)
}

func PermanentlyDeleteNovelHandler(c *gin.Context) {
	itemID := c.Param("itemId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	if err := service.PermanentlyDeleteNovel(itemID, userClaims.UserID); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Novel permanently deleted")
}
