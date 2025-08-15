package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// registerSettingRoutes 注册所有与设置模块相关的路由
func registerSettingRoutes(rg *gin.RouterGroup) {
	// MOCK API, not follow RESTful standard
	rg.GET("/api-providers", getApiProvidersHandler)
	rg.GET("/api-providers/modal", getModalProvidersHandler)
	rg.GET("/api-keys", getApiKeysHandler)
	rg.POST("/api-keys", addApiKeyHandler)
	rg.PUT("/api-keys/:id", updateApiKeyHandler)
	rg.DELETE("/api-keys/:id", deleteApiKeyHandler)
	rg.GET("/privacy/settings", getPrivacySettingsHandler)
	rg.PUT("/privacy/collection-settings", updatePrivacyCollectionSettingsHandler)
	rg.GET("/system/themes", getThemesHandler)
	rg.GET("/system/settings", getSystemSettingsHandler)
	rg.PATCH("/system/settings", saveSystemSettingHandler)
	rg.GET("/usage-logs", getUsageDataHandler)
	rg.GET("/users/settings", getUserSettingsHandler)

}

func getUserSettingsHandler(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       "user-123",
			"name":     "张三",
			"email":    "zhangsan@example.com",
			"avatar":   "https://i.pravatar.cc/150?u=a042581f4e29026704d",
			"plan":     "专业版",
			"phone":    "188-8888-8888",
			"region":   "中国大陆",
			"timezone": "Asia/Shanghai (UTC+8)",
			"bio":      "一个热爱写作的开发者。",
		},
		"notifications": []gin.H{
			{"id": 1, "title": "新评论与回复", "description": "当有人回复您的作品时通知我", "enabled": true},
			{"id": 2, "title": "系统公告", "description": "接收重要的系统更新和维护通知", "enabled": true},
			{"id": 3, "title": "活动与优惠", "description": "获取最新的平台活动和优惠信息", "enabled": false},
		},
		"securitySettings": []gin.H{
			{"title": "账户密码", "status": "上次更新于 3 个月前", "action": "修改"},
			{"title": "双重验证 (2FA)", "status": "未启用", "action": "启用"},
			{"title": "登录设备管理", "status": "当前有 2 个活跃会话", "action": "管理"},
			{"title": "登录历史", "status": "上次登录：2 小时前 (上海)", "action": "查看"},
		},
		"proPlanFeatures": []string{
			"无限次AI对话",
			"访问所有高级模型 (GPT-4, Claude 3)",
			"更快的响应速度",
			"优先技术支持",
		},
	})
}
func saveUserSettingsHandler(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getSystemSettingsHandler(c *gin.Context) {
	SettingStore.systemSettingsMutex.RLock()
	defer SettingStore.systemSettingsMutex.RUnlock()
	c.JSON(http.StatusOK, SettingStore.System)
}

func saveSystemSettingHandler(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getThemesHandler(c *gin.Context) {
	SettingStore.themesMutex.RLock()
	defer SettingStore.themesMutex.RUnlock()
	c.JSON(http.StatusOK, SettingStore.Themes)
}
func getPrivacySettingsHandler(c *gin.Context) {
	SettingStore.privacySettingsMutex.RLock()
	defer SettingStore.privacySettingsMutex.RUnlock()
	c.JSON(http.StatusOK, SettingStore.Privacy)
}

func updatePrivacyCollectionSettingsHandler(c *gin.Context) {
	// Mock implementation
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getApiProvidersHandler(c *gin.Context) {
	SettingStore.apiProvidersMutex.RLock()
	defer SettingStore.apiProvidersMutex.RUnlock()
	c.JSON(http.StatusOK, SettingStore.ApiProviders)
}

func getModalProvidersHandler(c *gin.Context) {
	SettingStore.modalProvidersMutex.RLock()
	defer SettingStore.modalProvidersMutex.RUnlock()
	c.JSON(http.StatusOK, SettingStore.ModalProviders)
}

func getApiKeysHandler(c *gin.Context) {
	SettingStore.apiKeysMutex.RLock()
	defer SettingStore.apiKeysMutex.RUnlock()
	keys := make([]ApiKey, 0, len(SettingStore.ApiKeys))
	for _, key := range SettingStore.ApiKeys {
		keys = append(keys, key)
	}
	c.JSON(http.StatusOK, keys)
}

func addApiKeyHandler(c *gin.Context) {
	var newKeyData struct {
		Provider      string `json:"provider"`
		ProviderShort string `json:"providerShort"`
		Name          string `json:"name"`
		Key           string `json:"key"`
		Model         string `json:"model"`
		BaseURL       string `json:"baseUrl"`
		Status        string `json:"status"`
	}
	if err := c.ShouldBindJSON(&newKeyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SettingStore.apiKeysMutex.Lock()
	defer SettingStore.apiKeysMutex.Unlock()

	newId := len(SettingStore.ApiKeys) + 10 // Simple ID generation for mock
	keyFragment := ""
	if len(newKeyData.Key) > 7 {
		keyFragment = newKeyData.Key[:5] + "••••" + newKeyData.Key[len(newKeyData.Key)-4:]
	}

	newApiKey := ApiKey{
		ID:            newId,
		Provider:      newKeyData.Provider,
		ProviderShort: newKeyData.ProviderShort,
		Name:          newKeyData.Name,
		KeyFragment:   keyFragment,
		Model:         newKeyData.Model,
		BaseURL:       newKeyData.BaseURL,
		Calls:         "0",
		Status:        newKeyData.Status,
		Created:       time.Now().Format("2006-01-02"),
	}

	SettingStore.ApiKeys[newId] = newApiKey
	c.JSON(http.StatusCreated, newApiKey)
}

func updateApiKeyHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var keyUpdateData ApiKey
	if err := c.ShouldBindJSON(&keyUpdateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SettingStore.apiKeysMutex.Lock()
	defer SettingStore.apiKeysMutex.Unlock()

	key, ok := SettingStore.ApiKeys[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}

	// Update fields from request
	key.Name = keyUpdateData.Name
	key.Model = keyUpdateData.Model
	key.Status = keyUpdateData.Status
	key.BaseURL = keyUpdateData.BaseURL

	// If a new key is provided, update the fragment
	if keyUpdateData.KeyFragment != "" && len(keyUpdateData.KeyFragment) > 7 {
		key.KeyFragment = keyUpdateData.KeyFragment[:5] + "••••" + keyUpdateData.KeyFragment[len(keyUpdateData.KeyFragment)-4:]
	}

	SettingStore.ApiKeys[id] = key
	c.JSON(http.StatusOK, key)
}

func deleteApiKeyHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	SettingStore.apiKeysMutex.Lock()
	defer SettingStore.apiKeysMutex.Unlock()
	deletedKey, ok := SettingStore.ApiKeys[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key not found"})
		return
	}
	delete(SettingStore.ApiKeys, id)
	c.JSON(http.StatusOK, deletedKey)
}
func getUsageDataHandler(c *gin.Context) {
	SettingStore.usageDataMutex.RLock()
	defer SettingStore.usageDataMutex.RUnlock()

	// Simulate pagination for logs
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	allLogs := SettingStore.Usage.Logs
	totalLogs := len(allLogs)
	totalPages := totalLogs / limit
	if totalLogs%limit != 0 {
		totalPages++
	}
	start := (page - 1) * limit
	end := start + limit
	if start > totalLogs {
		start = totalLogs
	}
	if end > totalLogs {
		end = totalLogs
	}
	paginatedLogs := allLogs[start:end]

	response := UsageData{
		Stats:      SettingStore.Usage.Stats,
		ChartData:  SettingStore.Usage.ChartData,
		Logs:       paginatedLogs,
		TotalLogs:  totalLogs,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}
