package main

import (
	"sync"
	"time"
)

// NovelDataStore 封装了小说模块所需的所有数据和锁
type NovelDataStore struct {
	Novels           map[string]NovelDashboardItem
	novelsMutex      sync.RWMutex
	Projects         map[string]NovelProject
	projectsMutex    sync.RWMutex
	Categories       []NovelCategory
	categoriesMutex  sync.RWMutex
	Trash            map[string]DeletedItem
	trashMutex       sync.RWMutex
	RecentItems      []RecentActivityItem
	recentItemsMutex sync.RWMutex
}

// initNovelData 初始化小说模块的模拟数据
func initNovelData() {
	categories := []NovelCategory{"科幻", "奇幻", "悬疑", "恐怖", "都市", "言情", "历史"}
	novels := map[string]NovelDashboardItem{
		"novel-1": {ID: "novel-1", Title: "星际漫游者", Description: "一部关于孤独宇航员在未知星系中寻找回家之路的科幻史诗。", Cover: "https://images.unsplash.com/photo-1589998059171-988d887df646?q=80&w=200", Status: NovelStatus{Text: "编辑中", Class: "bg-green-100 text-green-800"}, Tags: []NovelTag{{Text: "科幻", Class: "bg-blue-100 text-blue-800"}}, Chapters: 5, LastUpdated: "2小时前", Category: "科幻"},
		"novel-2": {ID: "novel-2", Title: "时间之沙", Description: "当历史可以被改写，一个历史学家必须阻止一个神秘组织抹去关键的历史事件。", Cover: "https://images.unsplash.com/photo-1544947950-fa07a98d237f?q=80&w=200", Status: NovelStatus{Text: "待审核", Class: "bg-yellow-100 text-yellow-800"}, Tags: []NovelTag{{Text: "悬疑", Class: "bg-purple-100 text-purple-800"}}, Chapters: 15, LastUpdated: "昨天", Category: "悬疑"},
	}

	starRangerCharacters := []map[string]interface{}{
		{"id": "char-calvin", "name": "卡尔文·里德", "avatar": "https://i.pravatar.cc/150?u=calvin", "identity": "主角, 探索者四号宇航员", "gender": "男", "age": 35, "faction": "人类联邦", "summary": "孤独的宇航员，在一次深空探索任务中遭遇意外，被迫独自在未知星系中寻找归途。性格坚毅、冷静，但内心深处对家园有着强烈的眷恋。", "notes": "设计灵感来源于电影《月球》和《星际穿越》。需要重点刻画其在长期孤独环境下的心理变化。", "status": "editing"},
		{"id": "char-aila", "name": "艾拉 (AILA)", "avatar": "", "identity": "AI, 飞船智能核心", "summary": "第五代通用人工智能，负责“探索者四号”的全部系统运作。逻辑至上，声音平稳无波澜。在与卡尔文的长期相处中，其程序底层开始出现不符合预期的、类似人类情感的逻辑萌芽。", "status": "completed"},
	}

	characterChildren := []*TreeNode{
		{ID: "characters-overview", Title: "角色总览", Type: "characters_overview", Icon: "fa-solid fa-users", Content: "", IsOverview: true, IsReadOnly: true},
	}
	for _, char := range starRangerCharacters {
		characterChildren = append(characterChildren, &TreeNode{
			ID: char["id"].(string), Title: char["name"].(string), Type: "character_item", Icon: "fa-regular fa-user",
			Content: "<h1>" + char["name"].(string) + "</h1><p>身份：" + char["identity"].(string) + "</p><p>简介：" + char["summary"].(string) + "</p>", OriginalData: char,
		})
	}

	projects := map[string]NovelProject{
		"novel-1": {
			Metadata: NovelMetadata{
				ID: "novel-1", Title: "星际漫游者", Description: "一部关于孤独宇航员在未知星系中寻找回家之路的科幻史诗。", Cover: novels["novel-1"].Cover,
				Tags: []NovelTag{{Text: "科幻", Class: "bg-blue-50 text-blue-700"}}, Status: "连载中", ReferenceNovelIds: []string{"ref-asimov-foundation", "ref-cixin-darkforest"},
			},
			DirectoryData: []Volume{
				{ID: "vol-1", Type: "volume", Title: "第一卷：星尘之始", Content: "<h1>第一卷：星尘之始</h1><p>本卷大纲...</p>", Chapters: []Chapter{
					{ID: "ch-1", Type: "chapter", Title: "第一章：深空孤影", WordCount: 3102, Content: "<h1>第一章：深空孤影</h1><p>这是章节的详细内容...</p>", Status: "completed"},
					{ID: "ch-2", Type: "chapter", Title: "第二章：异常信号", WordCount: 2845, Content: "<h1>第二章：异常信号</h1><p>一个神秘的信号打破了长久的平静...</p>", Status: "completed"},
					{ID: "ch-3", Type: "chapter", Title: "第三章：AI的低语", WordCount: 3500, Content: "<h1>第三章：AI的低语</h1><p>在分析信号的过程中，飞船的AI“艾拉”开始出现一些微小的异常行为...</p>", Status: "editing"},
					{ID: "ch-4", Type: "chapter", Title: "第四章: 跃迁点", WordCount: 2415, Content: "<h1>第四章: 跃迁点</h1><p>他们最终发现信号源自一个时空奇点...</p>", Status: "editing"},
				}},
				{ID: "vol-2", Type: "volume", Title: "第二卷：遗忘的航线", Content: "<h1>第二卷：遗忘的航线</h1><p>本卷大纲...</p>", Chapters: []Chapter{
					{ID: "ch-5", Type: "chapter", Title: "第五章：时空涟漪", WordCount: 0, Content: "<h1>第五章：时空涟漪</h1>", Status: "planned"},
				}},
			},
			SettingsData: []TreeNode{
				{ID: "setting", Title: "设定", Type: "root", Icon: "fa-solid fa-book-journal-whills", Children: []*TreeNode{
					{ID: "characters", Title: "角色", Type: "group", Icon: "fa-solid fa-users text-teal-500", Children: characterChildren},
					{ID: "locations", Title: "地点", Type: "group", Icon: "fa-solid fa-map-location-dot text-green-500", Children: []*TreeNode{{ID: "locations-overview", Title: "地点总览", Type: "locations_overview", Icon: "fa-solid fa-map-location-dot", IsOverview: true, IsReadOnly: true}}},
					{ID: "items", Title: "物品", Type: "group", Icon: "fa-solid fa-box-archive text-amber-600", Children: []*TreeNode{{ID: "items-overview", Title: "物品总览", Type: "items_overview", Icon: "fa-solid fa-box-archive", IsOverview: true, IsReadOnly: true}}},
					{ID: "worldview", Title: "世界观", Type: "group", Icon: "fa-solid fa-earth-americas text-sky-500", Children: []*TreeNode{{ID: "worldview-overview", Title: "世界观总览", Type: "worldview_overview", Icon: "fa-solid fa-book-atlas", IsOverview: true, IsReadOnly: true}, {ID: "world-overview-item", Title: "世界观细则", Type: "worldview_item", Icon: "fa-solid fa-book-atlas", Content: "<h1>世界观总览</h1><p>23世纪，人类掌握了亚光速航行技术...</p>"}}},
				}},
			},
			PlotCustomData:      []ItemNode{{ID: "custom-plot-1", Title: "关于跃迁点背后的文明猜想", Type: "plot_item", Icon: "fa-solid fa-lightbulb text-rose-500", Content: "<h1>关于跃迁点背后的文明猜想</h1>"}},
			AnalysisCustomData:  []ItemNode{},
			DerivedPlotData:     []PlotAnalysisItem{},
			DerivedAnalysisData: []PlotAnalysisItem{{ID: "analysis_1682495833181", Type: "analysis", SourceID: "ch-3", Title: "《第三章：AI的低语》分析", Content: "<h1>《第三章：AI的低语》分析</h1><p>本章通过AI“艾拉”的微小异常，成功地在科幻背景下引入了悬疑元素...</p>"}},
			OthersCustomData:    []ItemNode{{ID: "custom-others-1", Title: "写作风格参考", Type: "others_item", Icon: "fa-regular fa-file-zipper", Content: "<h1>写作风格参考</h1><p>参考阿西莫夫《基地》系列的宏大叙事风格...</p>"}},
			NoteData:            []NoteItem{{ID: "note-1", Type: "note", Title: "第四章情感转折点设计", Timestamp: "今天 14:32", Content: "<h1>第四章情感转折点设计</h1><p>需要重点描写卡尔文在面对跃迁点时，希望与恐惧交织的复杂心理。</p>"}},
		},
		"ref-asimov-foundation": {
			Metadata:      NovelMetadata{ID: "ref-asimov-foundation", Title: "《银河帝国：基地》", Description: "阿西莫夫的经典科幻作品", Cover: "https://images.unsplash.com/photo-1532012197267-da84d127e765?q=80&w=200", Tags: []NovelTag{{Text: "经典科幻", Class: "bg-gray-100 text-gray-800"}}, Status: "已完结"},
			DirectoryData: []Volume{{ID: "ref-asimov-vol-1", Type: "volume", Title: "第一部 心理史学家", Content: "<h1>第一部 心理史学家</h1>", Chapters: []Chapter{{ID: "ref-asimov-ch-1", Type: "chapter", Title: "第一节", WordCount: 3000, Content: "<h1>第一节</h1><p>哈里·谢顿在川陀的最后一次演讲...</p>", Status: "completed"}}}},
			SettingsData:  []TreeNode{}, PlotCustomData: []ItemNode{}, AnalysisCustomData: []ItemNode{}, DerivedPlotData: []PlotAnalysisItem{}, DerivedAnalysisData: []PlotAnalysisItem{}, OthersCustomData: []ItemNode{}, NoteData: []NoteItem{},
		},
		"ref-cixin-darkforest": {
			Metadata:      NovelMetadata{ID: "ref-cixin-darkforest", Title: "《黑暗森林》", Description: "刘慈欣《三体》系列的第二部", Cover: "https://images.unsplash.com/photo-1544716278-e513176f20b5?q=80&w=200", Tags: []NovelTag{{Text: "硬科幻", Class: "bg-blue-50 text-blue-700"}}, Status: "已完结"},
			DirectoryData: []Volume{{ID: "ref-cixin-vol-1", Type: "volume", Title: "面壁者", Content: "<h1>面壁者</h1>", Chapters: []Chapter{{ID: "ref-cixin-ch-1", Type: "chapter", Title: "第一章", WordCount: 5000, Content: "<h1>第一章</h1><p>面对三体文明的威胁，人类制定了“面壁计划”...</p>", Status: "completed"}}}},
			SettingsData:  []TreeNode{}, PlotCustomData: []ItemNode{}, AnalysisCustomData: []ItemNode{}, DerivedPlotData: []PlotAnalysisItem{}, DerivedAnalysisData: []PlotAnalysisItem{}, OthersCustomData: []ItemNode{}, NoteData: []NoteItem{},
		},
	}

	trashedItems := map[string]DeletedItem{
		"novel-3-deleted": {ID: "novel-3-deleted", Name: "深海回响", Type: "小说", Icon: "<svg fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" viewBox=\"0 0 24 24\"><rect x=\"5\" y=\"3\" width=\"14\" height=\"18\" rx=\"2\"></rect></svg>", DeletedAt: time.Now().Add(-48 * time.Hour).Format(time.RFC3339), RetentionDays: 27, RetentionPercent: 90},
	}
	recentItems := []RecentActivityItem{
		{ID: "activity-1", NovelID: "novel-1", NovelTitle: "星际漫游者", NovelCover: novels["novel-1"].Cover, EditedItemType: "chapter", EditedItemName: "第三章：AI的低语", EditedAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339), FormattedTime: "2小时前"},
		{ID: "activity-2", NovelID: "novel-2", NovelTitle: "时间之沙", NovelCover: novels["novel-2"].Cover, EditedItemType: "outline", EditedItemName: "第一卷大纲", EditedAt: time.Now().Add(-28 * time.Hour).Format(time.RFC3339), FormattedTime: "昨天"},
	}

	NovelStore = NovelDataStore{
		Novels:           novels,
		Projects:         projects,
		Categories:       categories,
		Trash:            trashedItems,
		RecentItems:      recentItems,
		novelsMutex:      sync.RWMutex{},
		projectsMutex:    sync.RWMutex{},
		categoriesMutex:  sync.RWMutex{},
		trashMutex:       sync.RWMutex{},
		recentItemsMutex: sync.RWMutex{},
	}
}
