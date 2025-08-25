package main

import (
	"fmt"
	"log"
	"st-novel-go/src/config"
	"st-novel-go/src/database"
	"st-novel-go/src/router"
)

// @title        AI Creator Platform API
// @version      1.0
// @description  This is the API server for the AI Creator Platform.

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

func main() {
	// Configuration is now loaded automatically via init() in the config package.

	// Initialize database connection
	database.InitDatabase()

	// Setup router
	r := router.SetupRouter()

	// Start server
	serverPort := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("Server is starting on port %s", serverPort)
	if err := r.Run(serverPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
