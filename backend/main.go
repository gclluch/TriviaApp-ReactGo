package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gclluch/TriviaApp-ReactGo/game"
	"github.com/gclluch/TriviaApp-ReactGo/handlers"
	"github.com/gclluch/TriviaApp-ReactGo/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load configurations
	loadConfig()

	// Initialize the Gin router with CORS configuration
	router := setupRouter()

	// Initialize the game server with preloaded questions
	gameServer := initializeGameServer()

	// Register HTTP and WebSocket handlers
	handlers.RegisterHandlers(router, gameServer)

	// Start the HTTP server
	startServer(router)
}

func loadConfig() {
	viper.SetDefault("PORT", "8080")
	// viper.SetDefault("QUESTIONS_FILE", "questions.json")
	viper.AutomaticEnv() // Read from environment variables
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	return router
}

func initializeGameServer() *game.GameServer {
	// questions, err := services.LoadQuestions(viper.GetString("QUESTIONS_FILE"))
	// if err != nil {
	// 	log.Fatalf("Failed to load questions: %v", err)
	// }
	sessionStore := store.NewSessionStore()
	return game.NewGameServer(sessionStore)
}

func startServer(router *gin.Engine) {
	port := viper.GetString("PORT")
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
