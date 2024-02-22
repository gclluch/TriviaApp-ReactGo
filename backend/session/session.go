package session

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// PlayerSession encapsulates the state and operations of a game session.
type PlayerSession struct {
	sync.Mutex
	Score             int                       // Single player score or multiplayer high score.
	Players           map[string]*models.Player // Players participating in the session.
	Connections       map[*websocket.Conn]bool  // Active WebSocket connections.
	Questions         []models.Question         // Session questions.
	AnsweredQuestions map[string]bool           // Tracks if a question has been answered correctly.
}

// NewPlayerSession initializes a new session with default values.
func NewPlayerSession() *PlayerSession {
	return &PlayerSession{
		Players:           make(map[string]*models.Player),
		Connections:       make(map[*websocket.Conn]bool),
		AnsweredQuestions: make(map[string]bool),
	}
}

// UpdateScore modifies the session's score and ensures thread safety.
func (ps *PlayerSession) UpdateScore(scoreToAdd int) {
	ps.Lock()
	defer ps.Unlock()

	ps.Score += scoreToAdd
}

// UpdatePlayerScore adjusts a player's score within a session.
func (ps *PlayerSession) UpdatePlayerScore(playerID, questionID string, scoreToAdd int) {
	ps.Lock()
	defer ps.Unlock()

	if player, exists := ps.Players[playerID]; exists {
		player.Score += scoreToAdd
		ps.AnsweredQuestions[questionID] = true // Mark the question as answered
	}
}

// AddPlayer introduces a new player to the session.
func (ps *PlayerSession) AddPlayer() *models.Player {
	ps.Lock()
	defer ps.Unlock()

	playerID := uuid.New().String()
	playerName := fmt.Sprintf("Player %d", len(ps.Players)+1)
	player := &models.Player{ID: playerID, Name: playerName}

	ps.Players[playerID] = player
	return player
}

// WEBSOCKET INTEGRATION

// AddConnection manages a new WebSocket connection.
func (ps *PlayerSession) AddConnection(conn *websocket.Conn) {
	ps.Lock()
	ps.Connections[conn] = true
	ps.Unlock()
	log.Println("New player connected.")
}

// Broadcast transmits messages to all active WebSocket connections.
func (ps *PlayerSession) Broadcast(message interface{}) {
	ps.Lock()
	defer ps.Unlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	for conn := range ps.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Failed to send message: %v", err)
			delete(ps.Connections, conn)
		}
	}
}

// BroadcastPlayerCount sends the current player count to all clients.
func (ps *PlayerSession) BroadcastPlayerCount() {
	message := map[string]interface{}{"type": "playerCount", "count": len(ps.Players)}
	ps.Broadcast(message)
}

// CheckAllPlayersFinished verifies if all players have completed the session.
func (ps *PlayerSession) CheckAllPlayersFinished() bool {
	ps.Lock()
	defer ps.Unlock()

	for _, player := range ps.Players {
		if !player.Finished {
			return false
		}
	}
	return true
}

// StartCountdown initiates a countdown, broadcasting updates to all clients.
func (ps *PlayerSession) StartCountdown(duration int) {
	for i := duration; i >= 0; i-- {
		ps.Broadcast(map[string]interface{}{"type": "countdown", "time": i})
		time.Sleep(time.Second)
	}
}

// BroadcastHighScore announces the session's high score to all clients.
func (ps *PlayerSession) BroadcastHighScore() {
	// ps.Lock()
	// defer ps.Unlock()

	highScore := 0
	for _, player := range ps.Players {
		if player.Score > highScore {
			highScore = player.Score
		}
	}
	ps.Broadcast(map[string]interface{}{"type": "highScore", "score": highScore})
}
