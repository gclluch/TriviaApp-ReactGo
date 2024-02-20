package store

import (
	"sync"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/services"
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
func (s *SessionStore) CreateSession(questions []models.Question, numQuestions int) string {
	s.Lock()
	defer s.Unlock()

	uniqueSessionID := generateSessionID()

	// Shuffle questions first to ensure randomness
	shuffledQuestions := services.ShuffleQuestions(questions)

	// Select the specified number of questions
	if numQuestions > len(shuffledQuestions) {
		numQuestions = len(shuffledQuestions)
	}
	selectedQuestions := shuffledQuestions[:numQuestions]

	playerSession := session.NewPlayerSession()
	playerSession.Questions = selectedQuestions

	s.Sessions[uniqueSessionID] = playerSession

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
