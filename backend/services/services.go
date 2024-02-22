package services

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/gclluch/captrivia_multiplayer/models"
)

// LoadQuestions loads questions from a JSON file specified by filename.
func LoadQuestions(filename string) ([]models.Question, error) {
	bytes, err := os.ReadFile(filename) // Use os.ReadFile
	if err != nil {
		return nil, err
	}

	var questions []models.Question
	if err := json.Unmarshal(bytes, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

// ShuffleQuestions randomizes the order of questions.
// Accepts a slice of Question structs and returns a new shuffled slice.
func ShuffleQuestions(questions []models.Question) []models.Question {
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
	return questions
}

// CheckAnswer verifies if the provided answer index for a question is correct.
// Returns a boolean indicating correctness and a boolean indicating if the question was found.
func CheckAnswer(questions []models.Question, questionID string, answerIndex int) (bool, bool) {
	for _, q := range questions {
		if q.ID == questionID {
			return q.CorrectIndex == answerIndex, true
		}
	}
	return false, false // Question not found
}
