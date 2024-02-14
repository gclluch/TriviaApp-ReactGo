package main

import (
	"log"
	"os"

	"github.com/gclluch/captrivia_multiplayer/game"
	"github.com/gclluch/captrivia_multiplayer/handlers"
	"github.com/gclluch/captrivia_multiplayer/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set up the Gin router and configure CORS, if needed.
	router := gin.Default()

	// Set up CORS middleware options
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Or configure as needed for your application

	// Apply CORS middleware to the router
	router.Use(cors.New(config))

	// Initialize the store.
	sessionStore := store.NewSessionStore()
	questions, err := game.LoadQuestions("questions.json")
	if err != nil {
		log.Fatalf("Failed to load questions: %v", err)
	}

	// Initialize the GameServer with its dependencies.
	gameServer := game.NewGameServer(sessionStore, questions)

	handlers.RegisterHandlers(router, gameServer)

	// Start the HTTP server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified.
	}

	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
