package mock

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// registerUserRoutes 注册所有与用户模块相关的路由
func registerUserRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", loginHandler)
		authGroup.POST("/register", registerHandler)
	}

	userGroup := rg.Group("/users")
	{
		userGroup.GET("/me", getCurrentUserHandler)
		userGroup.PUT("/settings", updateUserSettingsHandler)
	}
}

// loginHandler godoc
// @Summary User Login
// @Description Authenticates a user and returns a user object and a token.
// @Tags Users, Auth
// @Accept  json
// @Produce  json
// @Param   credentials body LoginCredentials true "Login Credentials"
// @Success 200 {object} map[string]interface{} "{"user": User, "token": "mock-token"}"
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func loginHandler(c *gin.Context) {
	var creds LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	UserStore.usersMutex.RLock()
	defer UserStore.usersMutex.RUnlock()

	// 模拟密码验证
	if user, ok := UserStore.Users["1"]; ok && user.Email == creds.Email && creds.Password == "123456" {
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": "mock-jwt-token-for-admin-123456",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
}

// registerHandler godoc
// @Summary User Registration
// @Description Registers a new user.
// @Tags Users, Auth
// @Accept  json
// @Produce  json
// @Param   registrationData body RegistrationData true "Registration Data"
// @Success 201 {object} map[string]interface{} "{"user": User, "token": "mock-token"}"
// @Router /auth/register [post]
func registerHandler(c *gin.Context) {
	var regData RegistrationData
	if err := c.ShouldBindJSON(&regData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 模拟注册，不实际存储
	newUser := User{
		ID:    "2",
		Name:  regData.FirstName + " " + regData.LastName,
		Email: regData.Email,
		Plan:  "免费版",
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  newUser,
		"token": "mock-jwt-token-for-" + newUser.Email,
	})
}

// getCurrentUserHandler godoc
// @Summary Get Current User
// @Description Get the currently authenticated user's details based on token.
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} User
// @Failure 404 {object} map[string]string
// @Router /users/me [get]
func getCurrentUserHandler(c *gin.Context) {
	UserStore.usersMutex.RLock()
	defer UserStore.usersMutex.RUnlock()

	user, ok := UserStore.Users["1"] // 模拟已登录用户
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// updateUserSettingsHandler godoc
// @Summary Update User Settings
// @Description Updates settings for the current user.
// @Tags Users
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   settings body UserSettings true "User Settings Object"
// @Success 200 {object} UserSettings
// @Failure 400 {object} map[string]string
// @Router /users/settings [put]
func updateUserSettingsHandler(c *gin.Context) {
	var updatedSettings UserSettings
	if err := c.ShouldBindJSON(&updatedSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	UserStore.userSettingsMutex.Lock()
	defer UserStore.userSettingsMutex.Unlock()

	// 更新存储中的数据
	UserStore.Settings = updatedSettings

	// 同时更新主用户列表中的信息，以保证数据一致性
	UserStore.usersMutex.Lock()
	if user, ok := UserStore.Users[updatedSettings.User.ID]; ok {
		user.Name = updatedSettings.User.Name
		user.Email = updatedSettings.User.Email
		user.Avatar = updatedSettings.User.Avatar
		user.Bio = updatedSettings.User.Bio
		user.Phone = updatedSettings.User.Phone
		user.Region = updatedSettings.User.Region
		user.Timezone = updatedSettings.User.Timezone
		UserStore.Users[updatedSettings.User.ID] = user
	}
	UserStore.usersMutex.Unlock()

	c.JSON(http.StatusOK, UserStore.Settings)
}
