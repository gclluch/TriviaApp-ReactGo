package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gclluch/captrivia_multiplayer/handlers"
	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var sessionStore *store.SessionStore // Global variable

func main() {
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
	router.GET("/join/:sessionId", gameServer.JoinGameHandler)

	// Register the WebSocket endpoint.
	router.GET("/ws", websocketEndpoint)

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all CORS origins
	},
}

func websocketEndpoint(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Deserialize the incoming WebSocket message
		var msg models.WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// Handle the message based on its action type
		switch msg.Action {
		case "join":
			handleJoinMessage(msg, conn)
			// Add other case handlers as needed
		}
	}
}

func handleJoinMessage(msg models.WebSocketMessage, conn *websocket.Conn) {
	session, exists := sessionStore.GetSession(msg.SessionID) // Replace sessionStore with sessionstore
	if !exists {
		log.Printf("Session not found: %s\n", msg.SessionID)
		// Optionally, send an error message back to the client
		return
	}

	playerID := msg.PlayerID
	if playerID == "" {
		playerID = uuid.New().String() // Generate a new ID if not provided
	}

	// Create a new Player instance with the WebSocket connection
	player := &models.Player{
		ID:    playerID,
		Conn:  conn,
		Score: 0, // Initialize score, adjust as necessary
	}

	// Add the player to the session
	session.AddPlayer(player)

	// Broadcast to all players in the session that a new player has joined
	session.Broadcast(models.WebSocketMessage{
		Action:    "update",
		SessionID: msg.SessionID,
		// Include other relevant information, such as the list of current players
	})
}
