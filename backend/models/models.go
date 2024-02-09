package models

// Question represents a trivia question, its options, and the index of the correct answer.
type Question struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

// PlayerSession represents a single player's game session, including their score.
type PlayerSession struct {
	ID    string `json:"id"` // This could be used for identifying the session
	Score int
}
