package store

import (
	"sync"

	"github.com/gclluch/captrivia_multiplayer/session"
	"github.com/google/uuid"
)

// SessionStore manages player sessions with thread-safe operations.
type SessionStore struct {
	sync.Mutex
	Sessions map[string]*session.PlayerSession
}

// NewSessionStore creates a new instance of SessionStore.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		Sessions: make(map[string]*session.PlayerSession),
	}
}

func generateSessionID() string {
	return uuid.New().String()
}

// CreateSession generates a new session with a unique ID and initializes the score.
func (s *SessionStore) CreateSession() string {
	s.Lock()
	defer s.Unlock()
	// Your session ID generation logic here
	uniqueSessionID := generateSessionID() // Enhance
	s.Sessions[uniqueSessionID] = &session.PlayerSession{Score: 0}

	// Need to initialize the shuffled list of questions ans store here
	// You might want to handle the case where the session already exists.
	// You might want to handle the case where the session creation fails.
	// You might want to handle the case where the session ID generation fails.
	// You might want to handle the case where the session ID is not unique.
	return uniqueSessionID
}

// GetSession retrieves a session by ID if it exists.
func (s *SessionStore) GetSession(sessionID string) (*session.PlayerSession, bool) {
	s.Lock()
	defer s.Unlock()
	session, exists := s.Sessions[sessionID]
	return session, exists
	// You might want to handle the case where the session doesn't exist.
	// You might want to handle the case where the session exists but is nil.
	// You might want to handle the case where the session retrieval fails.
	// You might want to handle the case where the session ID is invalid.
}

// UpdateSessionScore updates the score for a given session.
func (s *SessionStore) UpdateSessionScore(sessionID string, newScore int) {
	s.Lock()
	defer s.Unlock()

	if session, exists := s.Sessions[sessionID]; exists {
		session.Score = newScore
	}
	// You might want to handle the case where the session doesn't exist.
	// You might want to handle the case where the score update fails.
	// You might want to handle the case where the score is negative.
	// You might want to handle the case where the score is too high.
}
