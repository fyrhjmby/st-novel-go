// st-novel-go/src/novel/handler/related_content_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func getNovelJSONField(c *gin.Context, getterFunc func(string, uint) (interface{}, error)) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	data, err := getterFunc(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to get data: "+err.Error())
		return
	}
	utils.Success(c, data)
}

func GetSettingsDataHandler(c *gin.Context) {
	getNovelJSONField(c, service.GetSettingsData)
}

func GetPlotCustomDataHandler(c *gin.Context) {
	getNovelJSONField(c, service.GetPlotCustomData)
}

func GetAnalysisCustomDataHandler(c *gin.Context) {
	getNovelJSONField(c, service.GetAnalysisCustomData)
}

func GetOthersCustomDataHandler(c *gin.Context) {
	getNovelJSONField(c, service.GetOthersCustomData)
}

func updateNovelJSONField(c *gin.Context, fieldName string) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload []dto.TreeNodeDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, "Invalid JSON data: "+err.Error())
		return
	}

	err := service.UpdateNovelJSONField(novelID, userClaims.UserID, fieldName, payload)
	if err != nil {
		utils.Fail(c, "Failed to update data: "+err.Error())
		return
	}

	utils.Success(c, payload)
}

func UpdateSettingsDataHandler(c *gin.Context) {
	updateNovelJSONField(c, "settings_data")
}

func UpdatePlotCustomDataHandler(c *gin.Context) {
	updateNovelJSONField(c, "plot_custom_data")
}

func UpdateAnalysisCustomDataHandler(c *gin.Context) {
	updateNovelJSONField(c, "analysis_custom_data")
}

func UpdateOthersCustomDataHandler(c *gin.Context) {
	updateNovelJSONField(c, "others_custom_data")
}
