package store

import (
	"sync"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SessionStore manages player sessions with thread-safe operations.
type SessionStore struct {
	sync.Mutex
	Sessions map[string]*PlayerSession
}

// PlayerSession holds the state of a player's session including their score.
type PlayerSession struct {
	sync.Mutex
	Score       int
	Players     map[string]*models.Player // Keyed by player ID
	Connections map[*websocket.Conn]bool
}

// NewSessionStore creates a new instance of SessionStore.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		Sessions: make(map[string]*PlayerSession),
	}
}

func generateSessionID() string {
	// randBytes := make([]byte, 16)
	// rand.Read(randBytes)
	// return fmt.Sprintf("%x", randBytes)
	return uuid.New().String()
}

// CreateSession generates a new session with a unique ID and initializes the score.
func (s *SessionStore) CreateSession() string {
	s.Lock()
	defer s.Unlock()
	// Your session ID generation logic here
	uniqueSessionID := generateSessionID() // Enhance
	s.Sessions[uniqueSessionID] = &PlayerSession{Score: 0}
	return uniqueSessionID
}

// GetSession retrieves a session by ID if it exists.
func (s *SessionStore) GetSession(sessionID string) (*PlayerSession, bool) {
	s.Lock()
	defer s.Unlock()
	session, exists := s.Sessions[sessionID]
	return session, exists
}

// UpdateSessionScore updates the score for a given session.
func (s *SessionStore) UpdateSessionScore(sessionID string, newScore int) {
	s.Lock()
	defer s.Unlock()

	if session, exists := s.Sessions[sessionID]; exists {
		session.Score = newScore
	}
	// You might want to handle the case where the session doesn't exist.
}

// In your PlayerSession struct file
func (ps *PlayerSession) AddPlayer(player *models.Player) {
	ps.Lock() // Assuming your PlayerSession includes a sync.Mutex for concurrency control
	defer ps.Unlock()
	if ps.Players == nil {
		ps.Players = make(map[string]*models.Player)
	}
	ps.Players[player.ID] = player
}
