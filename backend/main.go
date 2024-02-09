package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gclluch/captrivia_multiplayer/db"
	"github.com/gclluch/captrivia_multiplayer/handlers"
	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize the database connection
	err := db.SetupDatabase() // Initialize the database connection
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the store.
	sessionStore := store.NewSessionStore()

	// Load or initialize your questions here. This is a placeholder.
	// You might load these from a JSON file, a database, or define them in code.
	questions, err := loadQuestions()
	if err != nil {
		log.Fatalf("Failed to load questions: %v", err)
	}

	// Initialize the GameServer with its dependencies.
	gameServer := handlers.NewGameServer(sessionStore, questions)

	// Set up the Gin router and configure CORS, if needed.
	router := gin.Default()

	// Set up CORS middleware options
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Or configure as needed for your application

	// Apply CORS middleware to the router
	router.Use(cors.New(config))

	// Define routes and associate them with handler functions.
	router.POST("/game/start", gameServer.StartGameHandler)
	router.GET("/questions", gameServer.QuestionsHandler)
	router.POST("/answer", gameServer.AnswerHandler)
	router.POST("/game/end", gameServer.EndGameHandler)

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

// loadQuestions should load your quiz questions.

func loadQuestions() ([]models.Question, error) {
	fileBytes, err := ioutil.ReadFile("questions.json")
	if err != nil {
		log.Fatalf("Unable to read questions file: %v", err)
		return nil, err
	}

	var questions []models.Question
	if err := json.Unmarshal(fileBytes, &questions); err != nil {
		log.Fatalf("Unable to unmarshal questions JSON: %v", err)
		return nil, err
	}

	return questions, nil
}
