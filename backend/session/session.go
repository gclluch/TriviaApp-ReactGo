package session

import (
	"sync"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gorilla/websocket"
)

// PlayerSession holds the state of a player's session including their score.
type PlayerSession struct {
	sync.Mutex
	Score       int
	Players     map[string]*models.Player // Keyed by player ID
	Connections map[*websocket.Conn]bool
	// You might want to add more fields here, such as a list of questions. (multi.single)
	// You might want to add more fields here, such as a game state.
}

// In your PlayerSession struct file
func (ps *PlayerSession) AddPlayer(player *models.Player) {
	ps.Lock() // Assuming your PlayerSession includes a sync.Mutex for concurrency control
	defer ps.Unlock()
	if ps.Players == nil {
		ps.Players = make(map[string]*models.Player)
	}
	ps.Players[player.ID] = player
	// You might want to handle the case where the player already exists.
	// You might want to handle the case where the player addition fails.
	// You might want to handle the case where the player ID is not unique.
}

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
