package mock

import "sync"

// UserDataStore 封装了用户模块所需的所有数据
type UserDataStore struct {
	Users             map[string]User
	usersMutex        sync.RWMutex
	Settings          UserSettings
	userSettingsMutex sync.RWMutex
}

// initUserData 初始化用户模块的模拟数据
func initUserData() {
	mockUser := User{
		ID:       "1",
		Name:     "张小明",
		Email:    "admin@example.com",
		Plan:     "专业版",
		Avatar:   "https://i.pravatar.cc/150?u=admin",
		Phone:    "+86 13800138000",
		Region:   "中国大陆",
		Timezone: "Asia/Shanghai",
		Bio:      "一名热爱科幻小说创作的开发者。",
	}

	mockUserSettings := UserSettings{
		User: mockUser,
		Notifications: []NotificationSetting{
			{ID: 1, Title: "系统更新通知", Description: "当平台有重要功能更新时通知我", Enabled: true},
			{ID: 2, Title: "AI任务完成提醒", Description: "当长时间运行的AI任务完成时发送通知", Enabled: true},
			{ID: 3, Title: "社区互动消息", Description: "有人回复或点赞我的分享时提醒我", Enabled: false},
		},
		SecuritySettings: []SecuritySetting{
			{Title: "两步验证 (2FA)", Status: "已启用", StatusClass: "text-green-600", Action: "管理"},
			{Title: "密码", Status: "上次更新于3个月前", StatusClass: "text-yellow-600", Action: "修改"},
		},
		ProPlanFeatures: []string{"无限云存储", "高级AI模型优先使用权", "导出无水印高清格式", "团队协作功能", "专属客服支持"},
	}

	UserStore = UserDataStore{
		Users:             map[string]User{"1": mockUser},
		Settings:          mockUserSettings,
		usersMutex:        sync.RWMutex{},
		userSettingsMutex: sync.RWMutex{},
	}
}
