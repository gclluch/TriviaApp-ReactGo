package models

// Question represents a single trivia question with multiple choice answers.
type Question struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

// PlayerSession holds the state of a player's session including their score.
type PlayerSession struct {
	Score int
}

// AnswerSubmission represents the payload for a player's answer submission.
type AnswerSubmission struct {
	SessionID  string `json:"sessionId"`
	QuestionID string `json:"questionId"`
	Answer     int    `json:"answer"`
}
