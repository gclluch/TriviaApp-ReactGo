package session

import (
	"fmt"
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

// func (ps *PlayerSession) AddPlayer(player *Player) {
// 	// Implementation remains similar, but also manage WebSocket connection
// }

// func (ps *PlayerSession) Broadcast(message interface{}) {
// 	ps.Lock()
// 	defer ps.Unlock()
// 	for conn := range ps.Connections {
// 		if err := conn.WriteJSON(message); err != nil {
// 			// Handle errors, possibly removing the connection
// 		}
// 	}
// }

// // New method to handle incoming WebSocket messages for this session
// func (ps *PlayerSession) HandleMessage(msg []byte, sender *websocket.Conn) {
// 	// Process the message and potentially broadcast updates
// }

// WEBSOCKET INTEGRATION

func (ps *PlayerSession) AddConnection(conn *websocket.Conn) {
	ps.Lock()
	defer ps.Unlock()

	if ps.Connections == nil {
		ps.Connections = make(map[*websocket.Conn]bool)
	}
	ps.Connections[conn] = true

	// Start a goroutine to listen for messages from this connection
	go ps.listen(conn)
}

func (ps *PlayerSession) listen(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		ps.removeConnection(conn)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // Connection closed or error occurred
		}

		// Handle the message, e.g., broadcast to other players in the session
		ps.handleMessage(msg, conn)
	}
}

func (ps *PlayerSession) handleMessage(msg []byte, sender *websocket.Conn) {
	// Logic to handle a new message received from a client
	// For example, broadcast it to all other connections in this session
}

func (ps *PlayerSession) Broadcast(message interface{}) {
	ps.Lock()
	defer ps.Unlock()

	// for conn := range ps.Connections {
	// 	// Send message to each connection
	// }
}

func (ps *PlayerSession) removeConnection(conn *websocket.Conn) {
	ps.Lock()
	defer ps.Unlock()

	delete(ps.Connections, conn)
}
