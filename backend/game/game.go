package game

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/services"
	"github.com/gclluch/captrivia_multiplayer/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GameServer encapsulates the game logic and data.
type GameServer struct {
	Store     *store.SessionStore
	Questions []models.Question
}

func NewGameServer(store *store.SessionStore, questions []models.Question) *GameServer {
	return &GameServer{
		Store:     store,
		Questions: questions,
	}
}

// StartGameHandler handles requests to start a new game session.
func (gs *GameServer) StartGameHandler(c *gin.Context) {
	sessionID := gs.Store.CreateSession()

	// Generate a shareable link. This could be as simple as appending the session ID
	// to a base URL. For real deployment, ensure your base URL matches your deployed frontend.
	baseURL := "http://localhost:3000/join/"
	shareableLink := baseURL + sessionID

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID, "shareableLink": shareableLink, "message": "Game started successfully."})

}

// QuestionsHandler returns a set of questions for the game.
func (gs *GameServer) QuestionsHandler(c *gin.Context) {
	shuffledQuestions := services.ShuffleQuestions(gs.Questions)
	// Limit the questions to a certain number if desired
	c.JSON(http.StatusOK, shuffledQuestions)
}

// AnswerHandler handles answer submissions and updates the player's score.
func (gs *GameServer) AnswerHandler(c *gin.Context) {
	var submission models.AnswerSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, exists := gs.Store.GetSession(submission.SessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Validate the answer and update the score
	correct, questionExists := services.CheckAnswer(gs.Questions, submission.QuestionID, submission.Answer)
	if !questionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if correct {
		session.Score += 10                                              // Assume each correct answer gives 10 points
		gs.Store.UpdateSessionScore(submission.SessionID, session.Score) // Implement this method in your SessionStore
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"currentScore": session.Score,
	})
}

// EndGameHandler concludes the game and returns the final score.
func (gs *GameServer) EndGameHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, exists := gs.Store.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Game ended successfully.",
		"finalScore": session.Score,
	})
}

// WEBSOCKET
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all CORS origins
	},
}

func (gs *GameServer) WebSocketEndpoint(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Listen for incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Process the message
		gs.handleWebSocketMessage(msg, conn)
	}
}

func (gs *GameServer) handleWebSocketMessage(msg []byte, conn *websocket.Conn) {
	// Deserialize the message into a structured format
	var message WebSocketMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Printf("Error unmarshalling WebSocket message: %v", err)
		return
	}

	// switch message.Action {
	// case "join":
	// 	gs.handleJoinGame(message, conn)
	// 	// Handle other actions such as "startCountdown", "submitAnswer", etc.
	// }
}

// Define WebSocketMessage struct in models.go
type WebSocketMessage struct {
	Action    string `json:"action"`
	SessionID string `json:"sessionId"`
	// Additional fields as needed
}
