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

// GameServer struct encapsulates dependencies for the game logic.
type GameServer struct {
	Store     *store.SessionStore
	Questions []models.Question
	Upgrader  websocket.Upgrader
}

// NewGameServer initializes a new GameServer instance.
func NewGameServer(store *store.SessionStore, questions []models.Question) *GameServer {
	return &GameServer{
		Store:     store,
		Questions: questions,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo purposes; adjust as necessary.
			},
		},
	}
}

// StartGameHandler initiates a new game session.
func (gs *GameServer) StartGameHandler(c *gin.Context) {
	var requestBody struct {
		NumQuestions int `json:"numQuestions"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		requestBody.NumQuestions = 10
	}

	numQuestions := services.Clamp(requestBody.NumQuestions, 1, len(gs.Questions))
	sessionID, err := gs.Store.CreateSession(gs.Questions, numQuestions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Generate a shareable link for the session.
	shareableLink := fmt.Sprintf("%s/join/%s", c.Request.Host, sessionID)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Game session created successfully.",
		"sessionId":     sessionID,
		"shareableLink": shareableLink,
	})
}

// JoinGameHandler adds a player to an existing game session.
func (gs *GameServer) JoinGameHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, ok := gs.retrieveSession(c, sessionID)
	if !ok {
		return
	}
	player := session.AddPlayer()

	// Broadcast the updated player count to all clients in the session
	session.BroadcastPlayerCount()

	// Start the countdown when the first player joins
	if len(session.Players) == 1 {
		go session.StartCountdown(5)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Player joined successfully.",
		"playerId":   player.ID,
		"playerName": player.Name,
	})
}

// QuestionsHandler returns a set of questions for the game.
func (gs *GameServer) QuestionsHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, ok := gs.retrieveSession(c, sessionID)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": session.Questions})
}

// AnswerHandler handles answer submissions and updates the player's score.
func (gs *GameServer) AnswerHandler(c *gin.Context) {
	var submission models.AnswerSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, ok := gs.retrieveSession(c, submission.SessionID)
	if !ok {
		return
	}

	// Validate the answer and update the score
	correct, questionExists := services.CheckAnswer(gs.Questions, submission.QuestionID, submission.Answer)
	if !questionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

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

	// If question answered correctly for the first time, update the player's score.
	if correct && !session.AnsweredQuestions[submission.QuestionID] {
		session.UpdatePlayerScore(
			submission.PlayerID,
			submission.QuestionID,
			10,
		)
		session.BroadcastHighScore()
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

	session, ok := gs.retrieveSession(c, requestBody.SessionID)
	if !ok {
		return
	}

	player, ok := session.Players[requestBody.PlayerID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}
	player.Finished = true

	if session.CheckAllPlayersFinished() {
		session.Broadcast(map[string]interface{}{"type": "sessionComplete"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player marked as finished"})
}

// EndGameHandler concludes the game and returns the final score.
func (gs *GameServer) EndGameHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, ok := gs.retrieveSession(c, sessionID)
	if !ok {
		return
	}

	// Clean up the session data

	c.JSON(http.StatusOK, gin.H{
		"message":    "Game ended successfully.",
		"finalScore": session.Score,
	})
}

// FinalScoresHandler handles the request for final scores of a session.
func (gs *GameServer) FinalScoresHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, ok := gs.retrieveSession(c, sessionID)
	if !ok {
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

	// fmt.Printf("Scores: %v, Winners: %v, High Score: %d\n", scores, winners, highScore)
	c.JSON(http.StatusOK, gin.H{
		"scores":    scores,
		"winners":   winners,
		"highScore": highScore,
	})
}

// retrieveSession retrieves a session by its unique ID
func (gs *GameServer) retrieveSession(c *gin.Context, sessionID string) (*session.PlayerSession, bool) {
	session, exists := gs.Store.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return nil, false
	}
	return session, true
}

// Websocket Integration

// WebSocketEndpoint upgrades an HTTP connection to a WebSocket connection and handles incoming WebSocket messages.
func (gs *GameServer) WebSocketEndpoint(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket protocol.
	conn, err := gs.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}()

	log.Println("WebSocket connection established")

	// Listen for messages on the WebSocket connection.
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}
		gs.processWebSocketMessage(msg, conn)
	}
}

// processWebSocketMessage unmarshals and processes a single WebSocket message.
func (gs *GameServer) processWebSocketMessage(msg []byte, conn *websocket.Conn) {
	var message map[string]interface{}
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Printf("Error unmarshalling WebSocket message: %v", err)
		return
	}

	if action, ok := message["action"].(string); ok {
		switch action {
		case "joinSession":
			gs.handleJoinSession(message, conn)
		default:
			log.Printf("Unhandled action type: %s", action)
		}
	} else {
		log.Println("WebSocket message does not contain an action type")
	}
}

// handleJoinSession processes a "joinSession" action from a WebSocket message.
func (gs *GameServer) handleJoinSession(message map[string]interface{}, conn *websocket.Conn) {
	sessionID, ok := message["sessionId"].(string)
	if !ok {
		log.Println("WebSocket joinSession message does not contain 'sessionId'")
		return
	}

	session, exists := gs.Store.GetSession(sessionID)
	if !exists {
		log.Printf("Session not found: %s", sessionID)
		return
	}

	session.AddConnection(conn)
	log.Printf("Player joined session: %s", sessionID)

	session.BroadcastPlayerCount()
}
