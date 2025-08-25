// st-novel-go/src/novel/handler/history_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func GetHistoryHandler(c *gin.Context) {
	documentID := c.Param("documentId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	history, err := service.GetHistory(documentID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to get history: "+err.Error())
		return
	}
	utils.Success(c, history)
}

func RestoreVersionHandler(c *gin.Context) {
	documentID := c.Param("documentId")
	versionID := c.Param("versionId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	err := service.RestoreVersion(documentID, versionID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to restore version: "+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Version restored successfully.")
}
