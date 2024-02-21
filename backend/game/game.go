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
	var requestBody struct {
		NumQuestions int `json:"numQuestions"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil || requestBody.NumQuestions <= 0 {
		requestBody.NumQuestions = 10
	}

	numQuestions := min(requestBody.NumQuestions, len(gs.Questions))

	sessionID := gs.Store.CreateSession(gs.Questions, numQuestions)

	baseURL := "http://localhost:3000/join/" // Read from environment
	shareableLink := baseURL + sessionID

	c.JSON(http.StatusOK, gin.H{
		"sessionId":     sessionID,
		"shareableLink": shareableLink,
		"message":       "Game started successfully.",
	})
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (gs *GameServer) JoinGameHandler(c *gin.Context) {
	session, ok := getSessionFromRequest(gs, c)
	if !ok {
		return
	}

	playerInfo := session.AddPlayer()

	// Start the countdown when the first player joins
	if len(session.Players) == 1 {
		go session.StartCountdown(5) // Start a 10-second countdown
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

	// Validate the answer and update the score
	correct, questionExists := services.CheckAnswer(gs.Questions, submission.QuestionID, submission.Answer)
	if !questionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// session.Lock()
	// defer session.Unlock()

	// Single Player logic
	if submission.PlayerID == "" {
		if correct {
			session.UpdateScore(10)
		}
		c.JSON(http.StatusOK, gin.H{"correct": correct, "currentScore": session.Score})
		return
	}

	// Multiplayer logic
	player, ok := session.Players[submission.PlayerID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	if correct && !session.AnsweredQuestions[submission.QuestionID] {
		session.UpdatePlayerScore(
			submission.PlayerID,
			submission.QuestionID,
			10,
		)
		session.BroadcastHighScore()
		// session.AnsweredQuestions[submission.QuestionID] = true
		// player.Score += 10 // Update individual player score
	}

	c.JSON(http.StatusOK, gin.H{"correct": correct, "currentScore": player.Score})
}

// MarkPlayerFinishedHandler updates a player's finished status and checks if all players are done.
func (gs *GameServer) MarkPlayerFinishedHandler(c *gin.Context) {
	var requestBody struct {
		SessionID string `json:"sessionId"`
		PlayerID  string `json:"playerId"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, exists := gs.Store.GetSession(requestBody.SessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	player, ok := session.Players[requestBody.PlayerID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	fmt.Println("Player finished: ", player.Name)

	player.Finished = true
	if session.CheckAllPlayersFinished() {
		fmt.Println("All players finished")
		session.Broadcast(map[string]interface{}{"type": "sessionComplete"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player marked as finished"})
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

// FinalScoresHandler handles the request for final scores of a session.
func (gs *GameServer) FinalScoresHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, exists := gs.Store.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Extract scores and determine winners
	var winners []string
	highScore := 0
	scores := make([]map[string]interface{}, 0)

	for _, player := range session.Players {
		scores = append(scores, map[string]interface{}{
			"playerName": player.Name,
			"score":      player.Score,
		})
		if player.Score > highScore {
			highScore = player.Score
			winners = []string{player.Name} // Reset winners with the new high scorer
		} else if player.Score == highScore {
			winners = append(winners, player.Name) // Add player to winners in case of a tie
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"scores":    scores,
		"winners":   winners,
		"highScore": highScore,
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
