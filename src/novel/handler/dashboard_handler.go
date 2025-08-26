package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func GetNovelsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	novels, err := service.GetNovels(userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to fetch novels")
		return
	}
	utils.Success(c, novels)
}

func CreateNovelHandler(c *gin.Context) {
	var payload dto.CreateNovelPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	newNovel, err := service.CreateNovel(payload, userClaims)
	if err != nil {
		utils.Fail(c, "Failed to create novel")
		return
	}
	utils.Success(c, newNovel)
}

func GetCategoriesHandler(c *gin.Context) {
	categories := service.GetAvailableCategories()
	utils.Success(c, categories)
}
