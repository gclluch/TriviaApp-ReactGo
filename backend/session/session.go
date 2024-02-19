package session

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// PlayerSession holds the state of a player's session including their score.
type PlayerSession struct {
	sync.Mutex
	Score       int
	Players     map[string]*models.Player // Keyed by player ID
	Connections map[*websocket.Conn]bool
	Questions   []models.Question // holds shuffled questions for session
	// You might want to add more fields here, such as a list of questions. (multi.single)
	// You might want to add more fields here, such as a game state.
}

// NewPlayerSession creates a new PlayerSession with initialized fields.
func NewPlayerSession() *PlayerSession {
	return &PlayerSession{
		Players:     make(map[string]*models.Player),
		Connections: make(map[*websocket.Conn]bool),
	}
}

// In your PlayerSession struct file
func (ps *PlayerSession) AddPlayer() models.Player {
	ps.Lock()
	defer ps.Unlock()
	if ps.Players == nil {
		ps.Players = make(map[string]*models.Player)
	}
	playerID := uuid.New().String()                     // Generate a unique ID for the player
	playerCount := len(ps.Players) + 1                  // Determine the player's number
	playerName := fmt.Sprintf("Player %d", playerCount) // Assign a name based on order

	player := &models.Player{ID: playerID, Name: playerName}
	ps.Players[playerID] = player

	return models.Player{ID: playerID, Name: playerName} // Return the new player's ID and name
}

// return playerID, player.Name
// You might want to handle the case where the player already exists.
// You might want to handle the case where the player addition fails.
// You might want to handle the case where the player ID is not unique.

// WEBSOCKET INTEGRATION

// AddConnection adds a new WebSocket connection to the session and starts listening for messages.
func (ps *PlayerSession) AddConnection(conn *websocket.Conn) {
	ps.Lock()
	defer ps.Unlock()

	if ps.Connections == nil {
		ps.Connections = make(map[*websocket.Conn]bool) // Initialize if nil
	}
	ps.Connections[conn] = true

	fmt.Println("Player connected to session")

	go ps.listen(conn) // Start listening for messages from this connection
}

// listen continuously reads messages from the WebSocket connection and handles them.
func (ps *PlayerSession) listen(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		ps.removeConnection(conn) // Ensure connection is removed on disconnect
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		ps.handleMessage(msg, conn)
	}
}

// handleMessage processes incoming WebSocket messages.
func (ps *PlayerSession) handleMessage(msg []byte, sender *websocket.Conn) {
	// Example: Log received messages
	fmt.Printf("Received message: %s\n", string(msg))

	// Here, add logic to handle different message types, e.g., update game state or broadcast messages
}

// Broadcast sends a message to all connected WebSocket clients in the session.
func (ps *PlayerSession) Broadcast(message interface{}) {
	ps.Lock()
	defer ps.Unlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	for conn := range ps.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Error broadcasting message: %v", err)
			delete(ps.Connections, conn) // Remove faulty connection
		}
	}
}

// BroadcastPlayerCount sends the current player count to all clients in the session.
func (ps *PlayerSession) BroadcastPlayerCount() {
	playerCount := len(ps.Players) // Determine the current player count
	message := map[string]interface{}{
		"type":  "playerCount",
		"count": playerCount,
	}
	ps.Broadcast(message) // Use the Broadcast method to send the message
}

// removeConnection safely removes a WebSocket connection from the session.
func (ps *PlayerSession) removeConnection(conn *websocket.Conn) {
	ps.Lock()
	delete(ps.Connections, conn)
	ps.Unlock()

	log.Println("Player disconnected from session")
}
