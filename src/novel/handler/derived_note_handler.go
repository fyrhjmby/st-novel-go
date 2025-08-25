// 文件: ..\st-novel-go\src\novel\handler\derived_note_handler.go
// Package handler contains the HTTP handlers for processing incoming requests.
package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"
)

// Derived Content Handlers

// GetDerivedContentHandler handles the request to fetch all derived content for a novel.
func GetDerivedContentHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	items, err := service.GetDerivedContentForNovel(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, items)
}

func CreateDerivedContentHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload dto.CreateDerivedContentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	item, err := service.CreateDerivedContent(userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, item)
}

func UpdateDerivedContentHandler(c *gin.Context) {
	itemID := c.Param("itemId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload dto.UpdateDerivedContentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	item, err := service.UpdateDerivedContent(itemID, userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, item)
}

func DeleteDerivedContentHandler(c *gin.Context) {
	itemID := c.Param("itemId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	err := service.DeleteDerivedContent(itemID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Derived content deleted successfully.")
}

// Note Handlers

// GetNotesHandler handles the request to fetch all notes for a novel.
func GetNotesHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	notes, err := service.GetNotesForNovel(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, notes)
}

func CreateNoteHandler(c *gin.Context) {
	novelID := c.Param("novelId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload dto.CreateNotePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	note, err := service.CreateNote(novelID, userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, note)
}

func UpdateNoteHandler(c *gin.Context) {
	noteID := c.Param("noteId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload dto.UpdateNotePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	note, err := service.UpdateNote(noteID, userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, note)
}

func DeleteNoteHandler(c *gin.Context) {
	noteID := c.Param("noteId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	err := service.DeleteNote(noteID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "Note deleted successfully.")
}
