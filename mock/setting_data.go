package mock

import "sync"

// SettingDataStore 封装了设置模块所需的所有数据和锁
type SettingDataStore struct {
	ApiProviders         []ApiProvider
	apiProvidersMutex    sync.RWMutex
	ApiKeys              map[int]ApiKey
	apiKeysMutex         sync.RWMutex
	ModalProviders       []ModalProvider
	modalProvidersMutex  sync.RWMutex
	System               SystemSettings
	systemSettingsMutex  sync.RWMutex
	Privacy              PrivacySettingsData
	privacySettingsMutex sync.RWMutex
	Themes               []Theme
	themesMutex          sync.RWMutex
	Usage                UsageData
	usageDataMutex       sync.RWMutex
}

// initSettingData 初始化设置模块的模拟数据
func initSettingData() {
	apiProviders := []ApiProvider{
		{Name: "OpenAI", ShortName: "OA", Description: "GPT系列模型", StatusText: "2个密钥", ActiveKeys: 2, TotalCalls: "15,677"},
		{Name: "Claude", ShortName: "CL", Description: "Anthropic AI", StatusText: "1个密钥", ActiveKeys: 1, TotalCalls: "8,912"},
		{Name: "Azure OpenAI", ShortName: "AZ", Description: "Microsoft Azure", StatusText: "未配置", ActiveKeys: 0, TotalCalls: "0"},
		{Name: "Google Gemini", ShortName: "GG", Description: "Google多模态模型", StatusText: "未配置", ActiveKeys: 0, TotalCalls: "0"},
	}
	apiKeys := map[int]ApiKey{
		1: {ID: 1, Provider: "OpenAI", ProviderShort: "OA", Name: "生产环境密钥 (GPT-4)", KeyFragment: "sk-••••1a2b", Model: "gpt-4-turbo", Calls: "12,456", Status: "启用", Created: "2024-03-01", BaseURL: "", Description: "用于生产环境的主要GPT-4密钥", Temperature: 0.7, MaxTokens: 4096},
		2: {ID: 2, Provider: "OpenAI", ProviderShort: "OA", Name: "测试环境密钥 (GPT-3.5)", KeyFragment: "sk-••••3c4d", Model: "gpt-3.5-turbo", Calls: "3,221", Status: "暂停", Created: "2024-05-15", BaseURL: "", Description: "用于内部测试的GPT-3.5密钥", Temperature: 0.9, MaxTokens: 2048},
		3: {ID: 3, Provider: "Claude", ProviderShort: "CL", Name: "主密钥 (Opus)", KeyFragment: "sk-ant-••••5e6f", Model: "claude-3-opus-20240229", Calls: "8,912", Status: "启用", Created: "2024-04-10", BaseURL: "", Description: "Claude 3 Opus模型的主力密钥", Temperature: 0.5, MaxTokens: 4096},
	}
	modalProviders := []ModalProvider{
		{Name: "OpenAI", ShortName: "OA", Description: "由OpenAI提供的行业领先模型，如GPT-4"},
		{Name: "Claude", ShortName: "CL", Description: "由Anthropic AI开发的大语言模型，擅长长文本处理"},
		{Name: "Azure OpenAI", ShortName: "AZ", Description: "企业级的OpenAI服务，具备更高安全性"},
		{Name: "Google Gemini", ShortName: "GG", Description: "Google出品的下一代多模态模型"},
	}
	themes := []Theme{
		{Name: "默认主题"},
		{Name: "护眼模式"},
		{Name: "深色模式"},
	}
	systemSettings := SystemSettings{
		ActiveTheme: "默认主题", ZoomLevel: 50, Language: "简体中文", DateFormat: "YYYY-MM-DD",
		NotificationSettings: []SettingItem{
			{Title: "启用桌面通知", Description: "在浏览器外接收重要提醒", Enabled: true},
			{Title: "启用邮件通知", Description: "将通知发送到您的注册邮箱", Enabled: false},
			{Title: "新功能与更新", Description: "接收关于平台新功能和更新的通知", Enabled: true},
		},
		AppSettings: []SettingItem{
			{Title: "自动保存", Description: "在您编辑时自动保存更改", Enabled: true},
			{Title: "拼写检查", Description: "在编辑器中启用实时拼写检查", Enabled: true},
			{Title: "对话历史", Description: "保存所有AI对话记录，方便回顾", Enabled: true},
		},
	}
	privacySettings := PrivacySettingsData{
		CollectionSettings: []DataCollectionSetting{
			{Title: "使用数据改进模型", Description: "允许我们使用您的非个人身份数据来改进我们的AI模型", Enabled: true},
			{Title: "收集使用诊断数据", Description: "自动发送崩溃报告和性能数据以帮助我们改进服务", Enabled: false},
			{Title: "个性化内容推荐", Description: "根据您的使用习惯推荐相关功能和模板", Enabled: true},
		},
		Usage: []DataUsageItem{
			{Title: "用户内容", Tag: "您的数据", Includes: "您创作的小说、笔记、设定等", Purpose: "为您提供核心的创作和存储服务"},
			{Title: "账户信息", Tag: "个人数据", Includes: "您的电子邮件、昵称、会员计划等", Purpose: "用于账户管理、认证和计费"},
		},
		Permissions: []DataPermission{
			{Title: "导出您的所有数据", Description: "随时可以打包并下载您的所有创作内容", Action: "立即导出"},
			{Title: "迁移您的数据", Description: "将您的数据迁移至其他兼容的服务", Action: "开始迁移"},
			{Title: "暂停数据收集", Description: "暂停所有非必要的数据收集活动", Action: "暂停收集"},
			{Title: "删除您的账户", Description: "永久删除您的账户和所有相关数据", Action: "请求删除"},
		},
		Promises: []string{"我们绝不会出售您的个人数据。", "我们仅在必要时访问您的数据以提供支持。", "您可以随时删除您的账户和所有数据。"},
	}
	usageData := UsageData{
		Stats: []UsageStat{
			{Label: "本月请求数", Value: "2,451", Trend: "+15%"},
			{Label: "本月Token消耗", Value: "1.2M", Trend: "+8%"},
			{Label: "平均响应时间", Value: "850ms", Trend: "-50ms"},
			{Label: "请求成功率", Value: "99.8%", Trend: "+0.1%"},
		},
		Logs: []ApiLog{
			{ID: 1, Timestamp: "2024-07-21 14:30:15", Endpoint: "/v1/chat/completions", Model: "gpt-4-turbo", Tokens: "1,204", Status: "成功", Duration: "1.2s"},
			{ID: 2, Timestamp: "2024-07-21 14:28:05", Endpoint: "/v1/messages", Model: "claude-3-opus", Tokens: "850", Status: "成功", Duration: "0.9s"},
			{ID: 3, Timestamp: "2024-07-21 14:27:01", Endpoint: "/v1/chat/completions", Model: "gpt-4-turbo", Tokens: "2,500", Status: "失败", Duration: "2.1s"},
		},
		ChartData: []ChartDataPoint{
			{Label: "7-15", Requests: 120, Tokens: 55}, {Label: "7-16", Requests: 150, Tokens: 75}, {Label: "7-17", Requests: 200, Tokens: 95},
			{Label: "7-18", Requests: 180, Tokens: 88}, {Label: "7-19", Requests: 220, Tokens: 110}, {Label: "7-20", Requests: 250, Tokens: 125},
			{Label: "7-21", Requests: 230, Tokens: 115},
		},
		TotalLogs: 25, TotalPages: 5,
	}

	SettingStore = SettingDataStore{
		ApiProviders:         apiProviders,
		ApiKeys:              apiKeys,
		ModalProviders:       modalProviders,
		System:               systemSettings,
		Privacy:              privacySettings,
		Themes:               themes,
		Usage:                usageData,
		apiProvidersMutex:    sync.RWMutex{},
		apiKeysMutex:         sync.RWMutex{},
		modalProvidersMutex:  sync.RWMutex{},
		systemSettingsMutex:  sync.RWMutex{},
		privacySettingsMutex: sync.RWMutex{},
		themesMutex:          sync.RWMutex{},
		usageDataMutex:       sync.RWMutex{},
	}
}
