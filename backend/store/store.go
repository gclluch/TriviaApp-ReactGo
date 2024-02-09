package store

import (
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/gclluch/captrivia_multiplayer/models"
)

// SessionStore manages player sessions with thread-safe operations.
type SessionStore struct {
	sync.Mutex
	Sessions map[string]*models.PlayerSession
}

// NewSessionStore creates a new instance of SessionStore.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		Sessions: make(map[string]*models.PlayerSession),
	}
}

func generateSessionID() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// CreateSession generates a new session with a unique ID and initializes the score.
func (s *SessionStore) CreateSession() string {
	s.Lock()
	defer s.Unlock()
	// Your session ID generation logic here
	uniqueSessionID := generateSessionID() // Enhance
	s.Sessions[uniqueSessionID] = &models.PlayerSession{Score: 0}
	return uniqueSessionID
}

// GetSession retrieves a session by ID if it exists.
func (s *SessionStore) GetSession(sessionID string) (*models.PlayerSession, bool) {
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
