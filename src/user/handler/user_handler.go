package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"st-novel-go/src/middleware"
	"st-novel-go/src/user/dao"
	"st-novel-go/src/user/service"
	"st-novel-go/src/utils"
)

func Register(c *gin.Context) {
	var payload service.RegisterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	user, err := service.Register(payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	user.Password = ""
	utils.Success(c, user)
}

func Login(c *gin.Context) {
	var payload service.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	response, err := service.Login(payload)
	if err != nil {
		utils.FailWithUnauthorized(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "login successful",
		"data": response,
	})
}

func GetCurrentUser(c *gin.Context) {
	claims, exists := c.Get(middleware.UserClaimsKey)
	if !exists {
		utils.FailWithUnauthorized(c, "Could not retrieve user claims")
		return
	}

	userClaims, ok := claims.(*utils.Claims)
	if !ok {
		utils.FailWithUnauthorized(c, "Invalid user claims format")
		return
	}

	user, err := dao.FindUserByEmail(userClaims.Email)
	if err != nil {
		utils.FailWithUnauthorized(c, "User not found")
		return
	}

	// Important: Do not expose password hash
	user.Password = ""
	utils.Success(c, user)
}

func GetUserSettingsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	settings, err := service.GetUserSettings(userClaims.UserID)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, settings)
}

func UpdateUserSettingsHandler(c *gin.Context) {
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var payload service.UpdateUserSettingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	err := service.UpdateUserSettings(userClaims.UserID, payload)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Settings updated successfully")
}
