package main

// --- 通用子类型 ---

// NovelStatus 定义小说状态及其前端样式
type NovelStatus struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

// NovelTag 定义小说标签及其前端样式
type NovelTag struct {
	Text  string `json:"text"`
	Class string `json:"class"`
}

// NovelCategory 定义小说分类，本质是字符串
type NovelCategory = string

// --- 仪表盘与列表 ---

// NovelDashboardItem 定义了在仪表盘上显示的小说条目
type NovelDashboardItem struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Cover       string        `json:"cover"`
	Status      NovelStatus   `json:"status"`
	Tags        []NovelTag    `json:"tags"`
	Chapters    int           `json:"chapters"`
	LastUpdated string        `json:"lastUpdated"`
	Category    NovelCategory `json:"category"`
	DeletedAt   string        `json:"deletedAt,omitempty"`
}

// DeletedItem 定义了回收站中的项目
type DeletedItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Icon             string `json:"icon"`
	DeletedAt        string `json:"deletedAt"`
	RetentionDays    int    `json:"retentionDays"`
	RetentionPercent int    `json:"retentionPercent"`
}

// RecentActivityItem 定义了最近活动列表中的项目
type RecentActivityItem struct {
	ID             string `json:"id"`
	NovelID        string `json:"novelId"`
	NovelTitle     string `json:"novelTitle"`
	NovelCover     string `json:"novelCover"`
	EditedItemType string `json:"editedItemType"`
	EditedItemName string `json:"editedItemName"`
	EditedAt       string `json:"editedAt"`
	FormattedTime  string `json:"formattedTime"`
}

// --- 小说项目编辑器内部结构 ---

// NovelMetadata 定义了小说的核心元数据
type NovelMetadata struct {
	ID                string     `json:"id"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	Cover             string     `json:"cover"`
	Tags              []NovelTag `json:"tags"`
	Status            string     `json:"status"`
	ReferenceNovelIds []string   `json:"referenceNovelIds"`
}

// Chapter 定义了章节
type Chapter struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // "chapter"
	Title     string `json:"title"`
	WordCount int    `json:"wordCount"`
	Content   string `json:"content"`
	Status    string `json:"status"`
}

// Volume 定义了分卷，包含多个章节
type Volume struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"` // "volume"
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Chapters []Chapter `json:"chapters"`
}

// NoteItem 定义了笔记
type NoteItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // "note"
	Title     string `json:"title"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// PlotAnalysisItem 定义了派生的大纲或分析项
type PlotAnalysisItem struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // "plot" or "analysis"
	SourceID string `json:"sourceId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// ItemNode 定义了自定义的、可自由组织的条目（如情节、分析、其他）
type ItemNode struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Type         string      `json:"type"`
	Icon         string      `json:"icon"`
	Content      string      `json:"content,omitempty"`
	OriginalData interface{} `json:"originalData,omitempty"`
}

// TreeNode 定义了编辑器左侧树形结构中的节点，用于设定等
type TreeNode struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Type         string      `json:"type"`
	Icon         string      `json:"icon"`
	Children     []*TreeNode `json:"children,omitempty"`
	Content      string      `json:"content,omitempty"`
	IsReadOnly   bool        `json:"isReadOnly,omitempty"`
	IsOverview   bool        `json:"isOverview,omitempty"`
	OriginalData interface{} `json:"originalData,omitempty"`
}

// NovelProject 定义了完整的小说项目，包含所有编辑数据
type NovelProject struct {
	Metadata            NovelMetadata      `json:"metadata"`
	DirectoryData       []Volume           `json:"directoryData"`
	SettingsData        []TreeNode         `json:"settingsData"`
	PlotCustomData      []ItemNode         `json:"plotCustomData"`
	AnalysisCustomData  []ItemNode         `json:"analysisCustomData"`
	DerivedPlotData     []PlotAnalysisItem `json:"derivedPlotData"`
	DerivedAnalysisData []PlotAnalysisItem `json:"derivedAnalysisData"`
	OthersCustomData    []ItemNode         `json:"othersCustomData"`
	NoteData            []NoteItem         `json:"noteData"`
}

// NovelProjectContent 定义了更新项目内容时使用的结构，不包含元数据
type NovelProjectContent struct {
	DirectoryData       []Volume           `json:"directoryData"`
	SettingsData        []TreeNode         `json:"settingsData"`
	PlotCustomData      []ItemNode         `json:"plotCustomData"`
	AnalysisCustomData  []ItemNode         `json:"analysisCustomData"`
	DerivedPlotData     []PlotAnalysisItem `json:"derivedPlotData"`
	DerivedAnalysisData []PlotAnalysisItem `json:"derivedAnalysisData"`
	OthersCustomData    []ItemNode         `json:"othersCustomData"`
	NoteData            []NoteItem         `json:"noteData"`
}
