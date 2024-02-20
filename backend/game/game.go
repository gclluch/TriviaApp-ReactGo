package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/services"
	"github.com/gclluch/captrivia_multiplayer/session"
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
	sessionID := gs.Store.CreateSession(gs.Questions)

	fmt.Println("Session ID: ", sessionID)

	// TODO: Pull this from env var or config
	baseURL := "http://localhost:3000/join/"
	shareableLink := baseURL + sessionID

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID, "shareableLink": shareableLink, "message": "Game started successfully."})
}

func (gs *GameServer) JoinGameHandler(c *gin.Context) {
	session, ok := getSessionFromRequest(gs, c)
	if !ok {
		return
	}

	playerInfo := session.AddPlayer()

	// Start the countdown when the first player joins
	if len(session.Players) == 1 {
		go session.StartCountdown(10) // Start a 10-second countdown
	}

	// Broadcast the updated player count to all clients in the session
	session.BroadcastPlayerCount() // Assuming this method broadcasts the player count

	c.JSON(http.StatusOK, gin.H{
		"message":    "Player joined successfully.",
		"playerId":   playerInfo.ID,   // Return the new player ID
		"playerName": playerInfo.Name, // Return the assigned player name
	})
}

// QuestionsHandler returns a set of questions for the game.
func (gs *GameServer) QuestionsHandler(c *gin.Context) {
	// TODO: Limit the questions to a certain number if desired
	session, ok := getSessionFromRequest(gs, c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, session.Questions)
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

	// fmt.Println("Answer submission: ", submission)
	// fmt.Println("Answer submission: ", session)

	// Validate the answer and update the score
	correct, questionExists := services.CheckAnswer(gs.Questions, submission.QuestionID, submission.Answer)
	if !questionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Single Player logic
	if submission.PlayerID == "" {
		if correct {
			session.Score += 10 // Assume each correct answer gives 10 points
		}
		c.JSON(http.StatusOK, gin.H{"correct": correct, "currentScore": session.Score})
		return
	}

	// Multiplayer logic
	session.Lock()
	defer session.Unlock()

	player, ok := session.Players[submission.PlayerID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	if correct && !session.AnsweredQuestions[submission.QuestionID] {
		session.AnsweredQuestions[submission.QuestionID] = true
		player.Score += 10           // Update individual player score
		session.BroadcastHighScore() // Implement this method to broadcast high score
	}

	c.JSON(http.StatusOK, gin.H{"correct": correct, "currentScore": player.Score})
}

// EndGameHandler concludes the game and returns the final score.
func (gs *GameServer) EndGameHandler(c *gin.Context) {
	session, ok := getSessionFromRequest(gs, c)
	if !ok {
		return
	}

	// Clean up the session data

	c.JSON(http.StatusOK, gin.H{
		"message":    "Game ended successfully.",
		"finalScore": session.Score,
	})
}

func getSessionFromRequest(gs *GameServer, c *gin.Context) (*session.PlayerSession, bool) {
	var sessionRequest models.SessionRequest
	if err := c.ShouldBindJSON(&sessionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return nil, false
	}

	session, exists := gs.Store.GetSession(sessionRequest.SessionId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return nil, false
	}

	return session, true
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

	// log.Println("WebSocket connection established")

	// Listen for incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s\n", message)

		// Correctly process the message here, within the loop
		gs.handleWebSocketMessage(message, conn)
	}
}

func (gs *GameServer) handleWebSocketMessage(msg []byte, conn *websocket.Conn) {
	// Assuming you have a way to parse your messages
	fmt.Println("Message received: ", string(msg))

	var message map[string]interface{}
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Printf("Error unmarshalling WebSocket message: %v", err)
		return
	}

	if action, ok := message["action"].(string); ok && action == "joinSession" {
		sessionId, _ := message["sessionId"].(string)
		// Retrieve the session using the sessionId
		session, exists := gs.Store.GetSession(sessionId)
		if !exists {
			log.Printf("Session not found: %s", sessionId)
			return
		}

		session.AddConnection(conn)
		fmt.Println("Player joined session: ", sessionId)

		// Optionally, send back the current player count
		playerCount := len(session.Players)
		conn.WriteJSON(map[string]interface{}{
			"type":  "playerCount",
			"count": playerCount,
		})
	}
}

func (gs *GameServer) BroadcastToSession(sessionId string, message interface{}) {
	session, exists := gs.Store.GetSession(sessionId)
	if !exists {
		log.Printf("Session not found: %s", sessionId)
		return
	}
	session.Broadcast(message)
}
