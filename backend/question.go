package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// Structs to parse the response
type APIResponse struct {
	Results []Question `json:"results"`
}

type Question struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

// Struct to format the question
type FormattedQuestion struct {
	ID           string   `json:"id"`
	QuestionText string   `json:"questionText"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
}

func main() {
	url := "https://opentdb.com/api.php?amount=10&type=multiple"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch questions:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Println("Failed to unmarshal JSON:", err)
		return
	}

	formattedQuestions := formatQuestions(apiResp.Results)
	saveToFile(formattedQuestions)
}

func formatQuestions(questions []Question) []FormattedQuestion {
	var formatted []FormattedQuestion
	for i, q := range questions {
		options := append(q.IncorrectAnswers, q.CorrectAnswer)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })

		// Finding the index of the correct answer after shuffling
		var correctIndex int
		for idx, opt := range options {
			if opt == q.CorrectAnswer {
				correctIndex = idx
				break
			}
		}

		formatted = append(formatted, FormattedQuestion{
			ID:           fmt.Sprintf("%d", i+1),
			QuestionText: q.Question,
			Options:      options,
			CorrectIndex: correctIndex,
		})
	}
	return formatted
}

func saveToFile(questions []FormattedQuestion) {
	file, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal questions:", err)
		return
	}

	if err := ioutil.WriteFile("triviaQuestions.json", file, 0644); err != nil {
		fmt.Println("Failed to write to file:", err)
	}
	fmt.Println("Trivia questions saved to triviaQuestions.json")
}
