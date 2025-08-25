package mock

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// registerNovelRoutes 注册所有与小说模块相关的路由
func registerNovelRoutes(rg *gin.RouterGroup) {
	novelGroup := rg.Group("/novels")
	{
		// Dashboard
		novelGroup.GET("", getNovelsHandler)
		novelGroup.POST("", createNovelHandler)
		novelGroup.GET("/categories", getNovelCategoriesHandler)
		novelGroup.DELETE("/:id", moveToTrashHandler)

		// Project Editor
		novelGroup.GET("/projects", getAllNovelProjectsHandler)
		novelGroup.GET("/:id/project", getNovelProjectHandler)
		novelGroup.PUT("/:id/project-content", updateNovelProjectContentHandler)
		novelGroup.PATCH("/:id/metadata", updateNovelMetadataHandler)
		novelGroup.POST("/import", importNovelProjectHandler)
		novelGroup.DELETE("/:id/permanent", deleteNovelProjectPermanentlyHandler)
	}

	// Trash
	trashGroup := rg.Group("/trash/novels")
	{
		trashGroup.GET("", getTrashedItemsHandler)
		trashGroup.POST("/:id/restore", restoreItemFromTrashHandler)
		trashGroup.DELETE("/:id", deleteItemFromTrashPermanentlyHandler)
	}

	// Recent Activities
	recentGroup := rg.Group("/recent-items")
	{
		recentGroup.GET("", getRecentItemsHandler)
		recentGroup.POST("", logRecentAccessHandler)
	}
}

// --- Dashboard Handlers ---

func getNovelsHandler(c *gin.Context) {
	NovelStore.novelsMutex.RLock()
	defer NovelStore.novelsMutex.RUnlock()
	novels := make([]NovelDashboardItem, 0, len(NovelStore.Novels))
	for _, n := range NovelStore.Novels {
		novels = append(novels, n)
	}
	c.JSON(http.StatusOK, novels)
}

func createNovelHandler(c *gin.Context) {
	var data struct {
		Title    string `json:"title"`
		Synopsis string `json:"synopsis"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	NovelStore.novelsMutex.Lock()
	NovelStore.projectsMutex.Lock()
	defer NovelStore.novelsMutex.Unlock()
	defer NovelStore.projectsMutex.Unlock()

	newID := "novel-" + uuid.New().String()
	newDashboardItem := NovelDashboardItem{
		ID: newID, Title: data.Title, Description: data.Synopsis, Category: data.Category,
		Cover:    "https://source.unsplash.com/random/200x280?book,sig=" + newID,
		Status:   NovelStatus{Text: "编辑中", Class: "bg-green-100 text-green-800"},
		Tags:     []NovelTag{{Text: data.Category, Class: "bg-blue-100 text-blue-800"}},
		Chapters: 1, LastUpdated: "刚刚",
	}

	newProject := NovelProject{
		Metadata: NovelMetadata{
			ID: newID, Title: data.Title, Description: data.Synopsis, Status: "连载中",
			Tags: newDashboardItem.Tags, Cover: newDashboardItem.Cover, ReferenceNovelIds: []string{},
		},
		DirectoryData: []Volume{
			{ID: "vol-" + uuid.New().String(), Type: "volume", Title: "第一卷", Content: "<h1>第一卷</h1>", Chapters: []Chapter{
				{ID: "ch-" + uuid.New().String(), Type: "chapter", Title: "第一章", WordCount: 0, Content: "<h1>第一章</h1>", Status: "editing"},
			}},
		},
		SettingsData: []TreeNode{}, PlotCustomData: []ItemNode{}, AnalysisCustomData: []ItemNode{},
		DerivedPlotData: []PlotAnalysisItem{}, DerivedAnalysisData: []PlotAnalysisItem{},
		OthersCustomData: []ItemNode{}, NoteData: []NoteItem{},
	}

	NovelStore.Novels[newID] = newDashboardItem
	NovelStore.Projects[newID] = newProject
	c.JSON(http.StatusCreated, newDashboardItem)
}

func getNovelCategoriesHandler(c *gin.Context) {
	NovelStore.categoriesMutex.RLock()
	defer NovelStore.categoriesMutex.RUnlock()
	c.JSON(http.StatusOK, NovelStore.Categories)
}

// --- Project Handlers ---

func getNovelProjectHandler(c *gin.Context) {
	id := c.Param("id")
	NovelStore.projectsMutex.RLock()
	defer NovelStore.projectsMutex.RUnlock()
	project, ok := NovelStore.Projects[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Novel project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func getAllNovelProjectsHandler(c *gin.Context) {
	NovelStore.projectsMutex.RLock()
	defer NovelStore.projectsMutex.RUnlock()
	projects := make([]NovelProject, 0, len(NovelStore.Projects))
	for _, p := range NovelStore.Projects {
		projects = append(projects, p)
	}
	c.JSON(http.StatusOK, projects)
}

func updateNovelProjectContentHandler(c *gin.Context) {
	id := c.Param("id")
	var contentData NovelProjectContent
	if err := c.ShouldBindJSON(&contentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	NovelStore.projectsMutex.Lock()
	defer NovelStore.projectsMutex.Unlock()

	project, ok := NovelStore.Projects[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	project.DirectoryData = contentData.DirectoryData
	project.SettingsData = contentData.SettingsData
	project.PlotCustomData = contentData.PlotCustomData
	project.AnalysisCustomData = contentData.AnalysisCustomData
	project.DerivedPlotData = contentData.DerivedPlotData
	project.DerivedAnalysisData = contentData.DerivedAnalysisData
	project.OthersCustomData = contentData.OthersCustomData
	project.NoteData = contentData.NoteData

	NovelStore.Projects[id] = project
	c.JSON(http.StatusOK, project)
}

func updateNovelMetadataHandler(c *gin.Context) {
	id := c.Param("id")
	var metadata NovelMetadata
	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	NovelStore.projectsMutex.Lock()
	NovelStore.novelsMutex.Lock()
	defer NovelStore.projectsMutex.Unlock()
	defer NovelStore.novelsMutex.Unlock()

	project, ok := NovelStore.Projects[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	project.Metadata = metadata
	NovelStore.Projects[id] = project

	if dashItem, ok := NovelStore.Novels[id]; ok {
		dashItem.Title = metadata.Title
		dashItem.Description = metadata.Description
		dashItem.Tags = metadata.Tags
		dashItem.Cover = metadata.Cover
		NovelStore.Novels[id] = dashItem
	}

	c.JSON(http.StatusOK, metadata)
}

func importNovelProjectHandler(c *gin.Context) {
	var importData struct {
		Title         string   `json:"title"`
		Description   string   `json:"description"`
		Category      string   `json:"category"`
		DirectoryData []Volume `json:"directoryData"`
	}

	if err := c.ShouldBindJSON(&importData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	NovelStore.novelsMutex.Lock()
	NovelStore.projectsMutex.Lock()
	defer NovelStore.novelsMutex.Unlock()
	defer NovelStore.projectsMutex.Unlock()

	newID := "novel-" + uuid.New().String()
	chapterCount := 0
	for _, vol := range importData.DirectoryData {
		chapterCount += len(vol.Chapters)
	}

	newDashboardItem := NovelDashboardItem{
		ID:          newID,
		Title:       importData.Title,
		Description: importData.Description,
		Category:    importData.Category,
		Cover:       "https://source.unsplash.com/random/200x280?book,sig=" + newID,
		Status:      NovelStatus{Text: "编辑中", Class: "bg-green-100 text-green-800"},
		Tags:        []NovelTag{{Text: importData.Category, Class: "bg-blue-100 text-blue-800"}},
		Chapters:    chapterCount,
		LastUpdated: "刚刚",
	}

	newProject := NovelProject{
		Metadata: NovelMetadata{
			ID:                newID,
			Title:             importData.Title,
			Description:       importData.Description,
			Status:            "连载中",
			Tags:              newDashboardItem.Tags,
			Cover:             newDashboardItem.Cover,
			ReferenceNovelIds: []string{},
		},
		DirectoryData:       importData.DirectoryData,
		SettingsData:        []TreeNode{}, // Start with empty settings
		PlotCustomData:      []ItemNode{},
		AnalysisCustomData:  []ItemNode{},
		DerivedPlotData:     []PlotAnalysisItem{},
		DerivedAnalysisData: []PlotAnalysisItem{},
		OthersCustomData:    []ItemNode{},
		NoteData:            []NoteItem{},
	}

	NovelStore.Novels[newID] = newDashboardItem
	NovelStore.Projects[newID] = newProject

	c.JSON(http.StatusCreated, newProject)
}

func deleteNovelProjectPermanentlyHandler(c *gin.Context) {
	id := c.Param("id")
	NovelStore.novelsMutex.Lock()
	NovelStore.projectsMutex.Lock()
	NovelStore.trashMutex.Lock()
	defer NovelStore.novelsMutex.Unlock()
	defer NovelStore.projectsMutex.Unlock()
	defer NovelStore.trashMutex.Unlock()

	delete(NovelStore.Novels, id)
	delete(NovelStore.Projects, id)
	delete(NovelStore.Trash, id)

	c.Status(http.StatusNoContent)
}

// --- Trash Handlers ---

func moveToTrashHandler(c *gin.Context) {
	id := c.Param("id")
	NovelStore.novelsMutex.Lock()
	NovelStore.trashMutex.Lock()
	defer NovelStore.novelsMutex.Unlock()
	defer NovelStore.trashMutex.Unlock()

	novel, ok := NovelStore.Novels[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
		return
	}
	delete(NovelStore.Novels, id)

	trashedItem := DeletedItem{
		ID: id, Name: novel.Title, Type: "小说",
		Icon:      "<svg fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" viewBox=\"0 0 24 24\"><rect x=\"5\" y=\"3\" width=\"14\" height=\"18\" rx=\"2\"/></svg>",
		DeletedAt: time.Now().Format(time.RFC3339), RetentionDays: 30, RetentionPercent: 100,
	}
	NovelStore.Trash[id] = trashedItem
	c.Status(http.StatusNoContent)
}

func getTrashedItemsHandler(c *gin.Context) {
	NovelStore.trashMutex.RLock()
	defer NovelStore.trashMutex.RUnlock()
	items := make([]DeletedItem, 0, len(NovelStore.Trash))
	for _, item := range NovelStore.Trash {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, items[i].DeletedAt)
		t2, _ := time.Parse(time.RFC3339, items[j].DeletedAt)
		return t1.After(t2)
	})
	c.JSON(http.StatusOK, items)
}

func restoreItemFromTrashHandler(c *gin.Context) {
	id := c.Param("id")
	NovelStore.novelsMutex.Lock()
	NovelStore.trashMutex.Lock()
	NovelStore.projectsMutex.RLock()
	defer NovelStore.novelsMutex.Unlock()
	defer NovelStore.trashMutex.Unlock()
	defer NovelStore.projectsMutex.RUnlock()

	trashedItem, ok := NovelStore.Trash[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not in trash"})
		return
	}
	delete(NovelStore.Trash, id)

	originalProject, ok := NovelStore.Projects[id]
	if !ok { // Fallback
		restoredNovel := NovelDashboardItem{ID: id, Title: trashedItem.Name, Description: "已恢复的小说。", LastUpdated: "刚刚", Status: NovelStatus{Text: "编辑中", Class: "bg-green-100 text-green-800"}}
		NovelStore.Novels[id] = restoredNovel
		c.JSON(http.StatusOK, restoredNovel)
		return
	}

	chapterCount := 0
	for _, vol := range originalProject.DirectoryData {
		chapterCount += len(vol.Chapters)
	}
	restoredDashboardItem := NovelDashboardItem{
		ID: id, Title: originalProject.Metadata.Title, Description: originalProject.Metadata.Description,
		Cover: originalProject.Metadata.Cover, Status: NovelStatus{Text: "编辑中", Class: "bg-green-100 text-green-800"},
		Tags: originalProject.Metadata.Tags, Chapters: chapterCount, LastUpdated: "刚刚", Category: "未分类",
	}
	NovelStore.Novels[id] = restoredDashboardItem
	c.JSON(http.StatusOK, restoredDashboardItem)
}

func deleteItemFromTrashPermanentlyHandler(c *gin.Context) {
	id := c.Param("id")
	NovelStore.trashMutex.Lock()
	NovelStore.projectsMutex.Lock()
	defer NovelStore.trashMutex.Unlock()
	defer NovelStore.projectsMutex.Unlock()

	delete(NovelStore.Trash, id)
	delete(NovelStore.Projects, id)
	c.Status(http.StatusNoContent)
}

// --- Recent Activities Handlers ---

func getRecentItemsHandler(c *gin.Context) {
	NovelStore.recentItemsMutex.RLock()
	defer NovelStore.recentItemsMutex.RUnlock()
	c.JSON(http.StatusOK, NovelStore.RecentItems)
}

func logRecentAccessHandler(c *gin.Context) {
	// Logic remains the same, just interacts with NovelStore
	// Implementation omitted for brevity
	c.JSON(http.StatusCreated, gin.H{"message": "Recent activity logged (mocked)"})
}
