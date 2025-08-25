// 文件: ..\st-novel-go\src\novel\router\novel_router.go

// st-novel-go/src/novel/router/novel_router.go
package router

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/handler"
)

func RegisterNovelRoutes(router *gin.RouterGroup) {
	novelRoutes := router.Group("/")
	novelRoutes.Use(middleware.AuthMiddleware())
	{
		// Dashboard & Trash routes
		novelsGroup := novelRoutes.Group("/novels")
		{
			novelsGroup.GET("", handler.GetNovelsHandler)
			novelsGroup.POST("", handler.CreateNovelHandler)
			novelsGroup.GET("/categories", handler.GetCategoriesHandler)
			novelsGroup.POST("/create-full", handler.CreateFullNovelProjectHandler) // New Route
			novelsGroup.DELETE("/:id", handler.MoveNovelToTrashHandler)
		}
		trashGroup := novelRoutes.Group("/trash")
		{
			trashGroup.GET("/novels", handler.GetTrashedNovelsHandler)
			trashGroup.POST("/novels/:id/restore", handler.RestoreNovelHandler)
			trashGroup.DELETE("/novels/:id", handler.PermanentlyDeleteNovelHandler)
		}

		// Import Route
		novelRoutes.POST("/novels/import", handler.ImportNovelProjectHandler)

		// Recent Activity routes
		recentGroup := novelRoutes.Group("/recent-items")
		{
			recentGroup.GET("", handler.GetRecentItemsHandler)
			recentGroup.POST("", handler.LogRecentAccessHandler)
		}

		// Novel Project routes
		projectGroup := novelRoutes.Group("/novels/projects")
		{
			projectGroup.GET("", handler.GetAllNovelProjectsHandler)
			projectGroup.GET("/:id", handler.GetNovelProjectHandler)
		}
		novelRoutes.DELETE("/novels/:id/permanent", handler.DeleteNovelProjectHandler)

		// Routes related to a specific Novel
		novelSpecificGroup := novelRoutes.Group("/novels/:novelId")
		{
			novelSpecificGroup.GET("/metadata", handler.GetNovelMetadataHandler)
			novelSpecificGroup.PATCH("/metadata", handler.UpdateNovelMetadataHandler)

			novelSpecificGroup.GET("/volumes", handler.GetVolumesHandler)
			novelSpecificGroup.POST("/volumes", handler.CreateVolumeHandler)
			novelSpecificGroup.PUT("/volumes/order", handler.UpdateVolumeOrderHandler)

			novelSpecificGroup.GET("/chapters", handler.GetChaptersForNovelHandler)

			novelSpecificGroup.GET("/settings", handler.GetSettingsDataHandler)
			novelSpecificGroup.PUT("/settings", handler.UpdateSettingsDataHandler)
			novelSpecificGroup.GET("/custom-plot", handler.GetPlotCustomDataHandler)
			novelSpecificGroup.PUT("/custom-plot", handler.UpdatePlotCustomDataHandler)
			novelSpecificGroup.GET("/custom-analysis", handler.GetAnalysisCustomDataHandler)
			novelSpecificGroup.PUT("/custom-analysis", handler.UpdateAnalysisCustomDataHandler)
			novelSpecificGroup.GET("/custom-others", handler.GetOthersCustomDataHandler)
			novelSpecificGroup.PUT("/custom-others", handler.UpdateOthersCustomDataHandler)

			novelSpecificGroup.GET("/derived-content", handler.GetDerivedContentHandler)
			novelSpecificGroup.GET("/notes", handler.GetNotesHandler)
			novelSpecificGroup.POST("/notes", handler.CreateNoteHandler)
		}

		// Routes for specific Volumes
		volumeSpecificGroup := novelRoutes.Group("/volumes/:volumeId")
		{
			volumeSpecificGroup.PATCH("", handler.UpdateVolumeHandler)
			volumeSpecificGroup.DELETE("", handler.DeleteVolumeHandler)
			volumeSpecificGroup.POST("/chapters", handler.CreateChapterHandler)
			volumeSpecificGroup.PUT("/chapters/order", handler.UpdateChapterOrderHandler)
		}

		// Routes for specific Chapters
		chapterSpecificGroup := novelRoutes.Group("/chapters/:chapterId")
		{
			chapterSpecificGroup.GET("", handler.GetChapterHandler)
			chapterSpecificGroup.PATCH("", handler.UpdateChapterHandler)
			chapterSpecificGroup.DELETE("", handler.DeleteChapterHandler)
		}

		// Routes for Derived Content and Notes
		derivedContentGroup := novelRoutes.Group("/derived-content")
		{
			derivedContentGroup.POST("", handler.CreateDerivedContentHandler)
			derivedContentGroup.PATCH("/:itemId", handler.UpdateDerivedContentHandler)
			derivedContentGroup.DELETE("/:itemId", handler.DeleteDerivedContentHandler)
		}
		notesGroup := novelRoutes.Group("/notes")
		{
			notesGroup.PATCH("/:noteId", handler.UpdateNoteHandler)
			notesGroup.DELETE("/:noteId", handler.DeleteNoteHandler)
		}

		// History routes
		documentsGroup := novelRoutes.Group("/documents/:documentId/history")
		{
			documentsGroup.GET("", handler.GetHistoryHandler)
			documentsGroup.POST("/:versionId/restore", handler.RestoreVersionHandler)
		}
	}
}
