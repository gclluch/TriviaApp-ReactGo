package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/gclluch/captrivia_multiplayer/models"
)

func LoadQuestions(filename string) ([]models.Question, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading questions file: %v", err)
		return nil, err
	}

	var questions []models.Question
	if err := json.Unmarshal(bytes, &questions); err != nil {
		log.Printf("Error unmarshalling questions: %v", err)
		return nil, err
	}

	return questions, nil
}

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
