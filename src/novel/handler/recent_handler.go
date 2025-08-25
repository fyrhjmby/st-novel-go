package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func GetRecentItemsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	items, err := service.GetRecentItems(userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to fetch recent items: "+err.Error())
		return
	}
	utils.Success(c, items)
}

func LogRecentAccessHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload dto.LogRecentAccessPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	activity, err := service.LogRecentAccess(payload, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to log recent access: "+err.Error())
		return
	}
	utils.Success(c, activity)
}
