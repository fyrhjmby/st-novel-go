package handler

import (
	"github.com/gin-gonic/gin"
	"st-novel-go/src/settings/service"
	"st-novel-go/src/utils"
)

// GetAPIProvidersHandler handles the request for the detailed provider list.
func GetAPIProvidersHandler(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	providers, err := service.GetAPIProviders(userID)
	if err != nil {
		utils.Fail(c, "Failed to get API providers")
		return
	}
	utils.Success(c, providers)
}

// GetModalProvidersHandler handles the request for the simplified provider list for modals.
func GetModalProvidersHandler(c *gin.Context) {
	providers := service.GetModalProviders()
	utils.Success(c, providers)
}
