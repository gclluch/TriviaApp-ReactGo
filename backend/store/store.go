package store

import (
	"fmt"
	"sync"

	"github.com/gclluch/TriviaApp-ReactGo/services"

	"github.com/gclluch/TriviaApp-ReactGo/session"
	"github.com/google/uuid"
)

// SessionStore manages player sessions and provides thread-safe operations to manipulate sessions.
type SessionStore struct {
	sync.Mutex
	Sessions map[string]*session.PlayerSession
}

// NewSessionStore initializes a new instance of SessionStore necessary defaults.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		Sessions: make(map[string]*session.PlayerSession),
	}
}

// CreateSession creates a new game session with a subset of questions and returns its unique ID.
// It shuffles the questions and selects the specified number to include in the session.
func (s *SessionStore) CreateSession(numQuestions int) (string, error) {
	if numQuestions <= 0 {
		return "", fmt.Errorf("numQuestions must be positive")
	}

	s.Lock()
	defer s.Unlock()

	// Generate a unique session ID.
	sessionID := uuid.New().String()

	questions, err := services.FetchQuestions(numQuestions)
	if err != nil {
		return "", err
	}

	playerSession := session.NewPlayerSession()

	playerSession.Questions = questions

	s.Sessions[sessionID] = playerSession

	fmt.Println(playerSession.Questions)
	return sessionID, nil
}

// GetSession retrieves a session by its unique ID, returning the session and a boolean indicating if it was found.
func (s *SessionStore) GetSession(sessionID string) (*session.PlayerSession, bool) {
	s.Lock()
	defer s.Unlock()

	session, exists := s.Sessions[sessionID]
	return session, exists
}
