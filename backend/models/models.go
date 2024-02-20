package models

// Question represents a single trivia question with multiple choice answers.
type Question struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

// AnswerSubmission represents the payload for a player's answer submission.
type AnswerSubmission struct {
	SessionID  string `json:"sessionId"`
	PlayerID   string `json:"playerId"`
	QuestionID string `json:"questionId"`
	Answer     int    `json:"answer"`
}

// Player represents a player in the game.
type Player struct {
	ID    string
	Name  string
	Score int
	// Conn  *websocket.Conn // WebSocket connection for real-time updates
	// Consider adding a WebSocket connection pointer here for direct messaging
}

// SessionRequest represents the payload for a request pertaining to a session
type SessionRequest struct {
	SessionId string `json:"sessionId"`
}
