package store

import (
	"fmt"
	"sync"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/services"
	"github.com/gclluch/captrivia_multiplayer/session"
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
func (s *SessionStore) CreateSession(questions []models.Question, numQuestions int) (string, error) {
	if numQuestions <= 0 {
		return "", fmt.Errorf("numQuestions must be positive")
	}

	s.Lock()
	defer s.Unlock()

	// Generate a unique session ID.
	sessionID := uuid.New().String()

	// Shuffle questions first to ensure randomness between sessions.
	shuffledQuestions := services.ShuffleQuestions(questions)

	// If numQuestions exceeds the total number of available questions, adjust it
	if numQuestions > len(shuffledQuestions) {
		numQuestions = len(shuffledQuestions)
	}

	// Select the desired number of questions for the session
	selectedQuestions := shuffledQuestions[:numQuestions]

	playerSession := session.NewPlayerSession()
	playerSession.Questions = selectedQuestions

	s.Sessions[sessionID] = playerSession

	return sessionID, nil
}

// GetSession retrieves a session by its unique ID, returning the session and a boolean indicating if it was found.
func (s *SessionStore) GetSession(sessionID string) (*session.PlayerSession, bool) {
	s.Lock()
	defer s.Unlock()

	session, exists := s.Sessions[sessionID]
	return session, exists
}
