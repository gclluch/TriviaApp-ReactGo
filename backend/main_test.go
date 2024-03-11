package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gclluch/TriviaApp-ReactGo/handlers"
	"github.com/gin-gonic/gin"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	loadConfig() // Assuming this sets up your environment as needed

	router := setupRouter()                       // Use the setup from your actual application
	gameServer := initializeGameServer()          // Initialize your game server with configurations
	handlers.RegisterHandlers(router, gameServer) // Register the same handlers as your application

	testServer = httptest.NewServer(router) // Use the configured router for testing

	exitCode := m.Run()

	testServer.Close()
	os.Exit(exitCode)
}

func TestStartGameHandler(t *testing.T) {
	// Simulating a POST request with JSON body
	body := strings.NewReader(`{"numQuestions": 10}`)
	resp, err := http.Post(fmt.Sprintf("%s/game/start", testServer.URL), "application/json", body)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check if a sessionId has been returned
	if _, exists := response["sessionId"]; !exists {
		t.Errorf("Response should contain a 'sessionId'")
	}
}

// func TestSinglePlayerGame(t *testing.T) {
// 	// Simulate a single player game

// 	// Step 1: Start a new single-player game
// 	startResp, err := http.Post(fmt.Sprintf("%s/game/start", testServer.URL), "application/json", nil)
// 	if err != nil {
// 		t.Fatalf("Failed to start a new game: %v", err)
// 	}
// 	defer startResp.Body.Close()

// 	if startResp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status OK; got %v", startResp.Status)
// 	}

// 	var gameSession map[string]interface{}
// 	if err := json.NewDecoder(startResp.Body).Decode(&gameSession); err != nil {
// 		t.Fatalf("Failed to decode start game response: %v", err)
// 	}

// 	sessionId, ok := gameSession["sessionId"].(string)
// 	if !ok {
// 		t.Fatalf("Session ID not found in start game response")
// 	}

// }
