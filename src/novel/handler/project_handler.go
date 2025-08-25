// 文件: ..\st-novel-go\src\novel\handler\project_handler.go

// st-novel-go/src/novel/handler/project_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

func GetNovelProjectHandler(c *gin.Context) {
	novelID := c.Param("id")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	project, err := service.GetNovelProject(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to get novel project: "+err.Error())
		return
	}
	utils.Success(c, project)
}

func GetAllNovelProjectsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	projects, err := service.GetAllNovelProjects(userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to get all novel projects: "+err.Error())
		return
	}
	utils.Success(c, projects)
}

func CreateFullNovelProjectHandler(c *gin.Context) {
	var payload dto.CreateFullNovelProjectPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	project, err := service.CreateFullNovelProject(payload, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to create full novel project: "+err.Error())
		return
	}
	utils.Success(c, project)
}

func ImportNovelProjectHandler(c *gin.Context) {
	var payload dto.CreateFullNovelProjectPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	project, err := service.ImportNovelProject(payload, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to import novel project: "+err.Error())
		return
	}
	utils.Success(c, project)
}

func DeleteNovelProjectHandler(c *gin.Context) {
	novelID := c.Param("id")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	err := service.PermanentlyDeleteNovel(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to permanently delete novel project: "+err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Novel project permanently deleted.")
}
