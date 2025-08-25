package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/settings/service"
	"st-novel-go/src/utils"
	"strconv"
)

// getUserID is a helper to safely extract user ID from the context.
func getUserID(c *gin.Context) (uint, bool) {
	claims, exists := c.Get(middleware.UserClaimsKey)
	if !exists {
		utils.FailWithUnauthorized(c, "Could not retrieve user claims")
		return 0, false
	}
	userClaims, ok := claims.(*utils.Claims)
	if !ok {
		utils.FailWithUnauthorized(c, "Invalid user claims format")
		return 0, false
	}
	return userClaims.UserID, true
}

func CreateAPIKeyHandler(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var payload service.CreateAPIKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	apiKey, err := service.CreateAPIKey(payload, userID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, apiKey)
}

func GetAPIKeysHandler(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	keys, err := service.GetAPIKeys(userID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, keys)
}

func UpdateAPIKeyHandler(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.FailWithBadRequest(c, "Invalid API key ID")
		return
	}

	var payload service.UpdateAPIKeyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	updatedKey, err := service.UpdateAPIKey(uint(id), userID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, updatedKey)
}

func DeleteAPIKeyHandler(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.FailWithBadRequest(c, "Invalid API key ID")
		return
	}

	_, err = service.DeleteAPIKey(uint(id), userID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "API key deleted successfully")
}
