package models

// Question represents a single trivia question with multiple choice answers.
type Question struct {
	ID           string   `json:"id"`           // Unique identifier for the question
	QuestionText string   `json:"questionText"` // The text of the question
	Options      []string `json:"options"`      // Available answers to the question
	CorrectIndex int      `json:"correctIndex"` // The index of the correct answer in the Options slice
}

// AnswerSubmission represents the payload for a player's answer submission.
type AnswerSubmission struct {
	SessionID  string `json:"sessionId"`  // Identifier for the game session
	PlayerID   string `json:"playerId"`   // Identifier for the player submitting the answer
	QuestionID string `json:"questionId"` // Identifier for the question being answered
	Answer     int    `json:"answer"`     // The index of the selected answer
}

// Player represents a player in the game.
type Player struct {
	ID       string `json:"id"`       // Unique identifier for the player
	Name     string `json:"name"`     // Name of the player
	Score    int    `json:"score"`    // Current score of the player
	Finished bool   `json:"finished"` // Whether the player has finished answering questions
}

// SessionRequest represents the payload for a request pertaining to a session.
type SessionRequest struct {
	SessionId string `json:"sessionId"` // Identifier for the session involved in the request
}

// Player Scores
type LeaderboardEntry struct {
	PlayerID       string `json:"playerId"`
	RightAnswers   int    `json:"RightAnswers"`
	TotalQuestions int    `json:"totalQuestions"`
}

// APIQuestion represents a question from the Open Trivia Database API.
type APIQuestion struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}
