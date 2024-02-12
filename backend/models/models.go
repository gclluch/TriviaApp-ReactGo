package models

import "sync"

// Question represents a single trivia question with multiple choice answers.
type Question struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

// Player represents a player in the game.
type Player struct {
	ID    string
	Score int
	// Consider adding a WebSocket connection pointer here for direct messaging
}

// AnswerSubmission represents the payload for a player's answer submission.
type AnswerSubmission struct {
	SessionID  string `json:"sessionId"`
	QuestionID string `json:"questionId"`
	Answer     int    `json:"answer"`
}

// PlayerSession holds the state of a player's session including their score.
type PlayerSession struct {
	sync.Mutex
	Score   int
	Players map[string]*Player // Keyed by player ID
	// Add more fields as needed (e.g., current question index)
}

// In your PlayerSession struct file
func (ps *PlayerSession) AddPlayer(player *Player) {
	ps.Lock() // Assuming your PlayerSession includes a sync.Mutex for concurrency control
	defer ps.Unlock()
	if ps.Players == nil {
		ps.Players = make(map[string]*Player)
	}
	ps.Players[player.ID] = player
}
