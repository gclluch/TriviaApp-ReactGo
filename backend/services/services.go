package services

import (
	"math/rand"
	"time"

	"github.com/gclluch/captrivia_multiplayer/models"
)

// ShuffleQuestions shuffles the slice of questions to randomize their order.
func ShuffleQuestions(questions []models.Question) []models.Question {
	rand.Seed(time.Now().UnixNano())
	shuffled := make([]models.Question, len(questions))
	copy(shuffled, questions)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

// CheckAnswer checks if the given answer for a question is correct.
// Returns true if correct, and also indicates if the question exists.
func CheckAnswer(questions []models.Question, questionID string, answerIndex int) (bool, bool) {
	for _, q := range questions {
		if q.ID == questionID {
			return q.CorrectIndex == answerIndex, true
		}
	}
	return false, false // Question not found
}

// UpdateSessionScore could be defined here if you want to include
// more complex logic around score calculation. However, it might
// involve changing how you interact with the session store,
// so it's not included in this basic example.
