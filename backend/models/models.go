package models

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

// PlayerSession holds the state of a player's session including their score.
type PlayerSession struct {
	Score   int
	Players map[string]*Player // Keyed by player ID
	// Add more fields as needed (e.g., current question index)
}

// AnswerSubmission represents the payload for a player's answer submission.
type AnswerSubmission struct {
	SessionID  string `json:"sessionId"`
	QuestionID string `json:"questionId"`
	Answer     int    `json:"answer"`
}
