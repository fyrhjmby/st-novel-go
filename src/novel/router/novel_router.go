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
		// Dashboard routes
		dashboardGroup := novelRoutes.Group("/novels")
		{
			dashboardGroup.GET("", handler.GetNovelsHandler)
			dashboardGroup.POST("", handler.CreateNovelHandler)
			dashboardGroup.GET("/categories", handler.GetCategoriesHandler)
		}

		// Trash routes
		trashGroup := novelRoutes.Group("/trash/novels")
		{
			trashGroup.GET("", handler.GetTrashedNovelsHandler)
			trashGroup.POST("/:itemId/restore", handler.RestoreNovelHandler)
			trashGroup.DELETE("/:itemId", handler.PermanentlyDeleteNovelHandler)
		}

		// Recent Activity routes
		recentGroup := novelRoutes.Group("/recent-items")
		{
			recentGroup.GET("", handler.GetRecentItemsHandler)
			recentGroup.POST("", handler.LogRecentAccessHandler)
		}

		// Novel Project and Import/Export routes
		novelRoutes.POST("/novels/import", handler.ImportNovelProjectHandler)
		novelRoutes.POST("/novels/create-full", handler.CreateFullNovelProjectHandler)
		projectGroup := novelRoutes.Group("/novels/projects")
		{
			projectGroup.GET("", handler.GetAllNovelProjectsHandler)
			projectGroup.GET("/:novelId", handler.GetNovelProjectHandler)
		}
		novelRoutes.DELETE("/novels/:novelId/permanent", handler.DeleteNovelProjectHandler)
		novelRoutes.DELETE("/novels/:novelId", handler.MoveNovelToTrashHandler)

		// Routes related to a specific Novel (Settings, Custom Data, etc.)
		novelSpecificGroup := novelRoutes.Group("/novels/:novelId")
		{
			novelSpecificGroup.GET("/metadata", handler.GetNovelMetadataHandler)
			novelSpecificGroup.PATCH("/metadata", handler.UpdateNovelMetadataHandler)
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
			novelSpecificGroup.GET("/volumes", handler.GetVolumesHandler)
			novelSpecificGroup.POST("/volumes", handler.CreateVolumeHandler)
			novelSpecificGroup.PUT("/volumes/order", handler.UpdateVolumeOrderHandler)
		}

		// Routes for specific Volumes
		volumeSpecificGroup := novelRoutes.Group("/volumes/:volumeId")
		{
			volumeSpecificGroup.PATCH("", handler.UpdateVolumeHandler)
			volumeSpecificGroup.DELETE("", handler.DeleteVolumeHandler)
			volumeSpecificGroup.GET("/chapters", handler.GetChaptersForVolumeHandler)
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

		// Routes for Derived Content and Notes (CRUD on individual items)
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
