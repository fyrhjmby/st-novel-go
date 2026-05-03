// stub_handlers.go — 系统设置、使用日志、数据隐私等 stub 端点
package handler

import (
	"st-novel-go/src/utils"

	"github.com/gin-gonic/gin"
)

// --- System Settings Stubs ---

func GetSystemThemesHandler(c *gin.Context) {
	themes := []map[string]interface{}{
		{"id": "default", "name": "默认", "type": "light"},
		{"id": "dark", "name": "深色", "type": "dark"},
		{"id": "eyecare-green", "name": "护眼绿", "type": "light"},
	}
	utils.Success(c, themes)
}

func GetSystemSettingsHandler(c *gin.Context) {
	settings := map[string]interface{}{
		"language":   "zh-CN",
		"dateFormat": "YYYY-MM-DD",
		"timezone":   "Asia/Shanghai",
		"zoomLevel":  100,
	}
	utils.Success(c, settings)
}

func UpdateSystemSettingsHandler(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, "Invalid settings data")
		return
	}
	utils.Success(c, payload)
}

// --- Usage Logs Stubs ---

func GetUsageLogsHandler(c *gin.Context) {
	logs := []map[string]interface{}{
		{"id": "1", "action": "AI任务-续写", "timestamp": "2026-05-03T10:00:00Z", "details": "续写第3章"},
		{"id": "2", "action": "导出", "timestamp": "2026-05-02T15:30:00Z", "details": "导出为TXT"},
	}
	utils.Success(c, logs)
}

// --- Data Privacy Stubs ---

func GetPrivacySettingsHandler(c *gin.Context) {
	settings := map[string]interface{}{
		"dataCollection":  true,
		"usageAnalytics":  true,
		"crashReports":    true,
		"personalization": false,
	}
	utils.Success(c, settings)
}

func UpdatePrivacySettingsHandler(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.FailWithBadRequest(c, "Invalid privacy settings")
		return
	}
	utils.Success(c, payload)
}
