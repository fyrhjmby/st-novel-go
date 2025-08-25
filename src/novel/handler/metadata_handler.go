// ..\st-novel-go\src\novel\handler\metadata_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func GetNovelMetadataHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	metadata, err := service.GetNovelMetadata(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, metadata)
}

func UpdateNovelMetadataHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload dto.UpdateMetadataPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	updatedMetadata, err := service.UpdateNovelMetadata(novelID, userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, updatedMetadata)
}
