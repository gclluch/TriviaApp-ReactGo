package services

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

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

// Helper function to get the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to clamp value between min and max.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func FetchQuestions(amount int) ([]models.Question, error) {
	apiQuestions, err := fetchAPIQuestions(amount)
	if err != nil {
		return nil, err
	}

	return FormatQuestions(apiQuestions), nil
}

// fetchAPIQuestions fetches trivia questions from an external API.
func fetchAPIQuestions(amount int) ([]models.APIQuestion, error) {
	var apiResponse struct {
		Results []models.APIQuestion `json:"results"`
	}

	url := fmt.Sprintf("https://opentdb.com/api.php?amount=%d&type=multiple", amount)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.Results, nil
}

// FormatQuestions formats a slice of APIQuestion into a slice of Question.
func FormatQuestions(apiQuestions []models.APIQuestion) []models.Question {
	var questions []models.Question

	// Initialize the random seed outside the loop
	rand.Seed(time.Now().UnixNano())

	for i, apiQ := range apiQuestions {
		// Decode HTML entities in question text
		questionText := html.UnescapeString(apiQ.Question)

		// Prepare options and decode HTML entities
		options := make([]string, len(apiQ.IncorrectAnswers)+1)
		for j, opt := range apiQ.IncorrectAnswers {
			options[j] = html.UnescapeString(opt)
		}
		correctOption := html.UnescapeString(apiQ.CorrectAnswer)
		options[len(options)-1] = correctOption

		// Shuffle options
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })

		// Find the index of the correct answer after shuffling
		correctIndex := findCorrectIndex(options, correctOption)

		questions = append(questions, models.Question{
			ID:           fmt.Sprintf("%d", i+1),
			QuestionText: questionText,
			Options:      options,
			CorrectIndex: correctIndex,
		})
	}

	return questions
}

// findCorrectIndex finds the index of the correct answer in the shuffled options.
func findCorrectIndex(options []string, correctAnswer string) int {
	for i, option := range options {
		if option == correctAnswer {
			return i
		}
	}
	return -1 // Indicates an error in finding the correct answer
}
