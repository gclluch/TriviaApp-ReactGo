package handlers

import (
	"net/http"

	"github.com/gclluch/captrivia_multiplayer/models"
	"github.com/gclluch/captrivia_multiplayer/services"
	"github.com/gclluch/captrivia_multiplayer/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GameServer encapsulates the game logic and data.
type GameServer struct {
	Store     *store.SessionStore
	Questions []models.Question // Assuming you have a way to load these questions
}

// NewGameServer creates a new GameServer instance.
func NewGameServer(store *store.SessionStore, questions []models.Question) *GameServer {
	return &GameServer{
		Store:     store,
		Questions: questions,
	}
}

// StartGameHandler handles requests to start a new game session.
func (gs *GameServer) StartGameHandler(c *gin.Context) {
	sessionID := gs.Store.CreateSession()

	// Generate a shareable link. This could be as simple as appending the session ID
	// to a base URL. For real deployment, ensure your base URL matches your deployed frontend.
	baseURL := "http://localhost:3000/join/"
	shareableLink := baseURL + sessionID

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID, "shareableLink": shareableLink, "message": "Game started successfully."})

	// c.JSON(http.StatusOK, gin.H{"sessionId": sessionID, "message": "Game started successfully."})
}

// QuestionsHandler returns a set of questions for the game.
func (gs *GameServer) QuestionsHandler(c *gin.Context) {
	shuffledQuestions := services.ShuffleQuestions(gs.Questions)
	// Limit the questions to a certain number if desired
	c.JSON(http.StatusOK, shuffledQuestions)
}

// AnswerHandler handles answer submissions and updates the player's score.
func (gs *GameServer) AnswerHandler(c *gin.Context) {
	var submission models.AnswerSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, exists := gs.Store.GetSession(submission.SessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Validate the answer and update the score
	correct, questionExists := services.CheckAnswer(gs.Questions, submission.QuestionID, submission.Answer)
	if !questionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if correct {
		session.Score += 10                                              // Assume each correct answer gives 10 points
		gs.Store.UpdateSessionScore(submission.SessionID, session.Score) // Implement this method in your SessionStore
	}

	c.JSON(http.StatusOK, gin.H{
		"correct":      correct,
		"currentScore": session.Score,
	})
}

// EndGameHandler concludes the game and returns the final score.
func (gs *GameServer) EndGameHandler(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, exists := gs.Store.GetSession(sessionID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Game ended successfully.",
		"finalScore": session.Score,
	})
}

// JoinSessionHandler handles requests to join an existing game session.
func (gs *GameServer) JoinGameHandler(c *gin.Context) {
	sessionId := c.Param("sessionId")

	playerID := uuid.New().String() // Implement this function based on your ID generation strategy

	player := &models.Player{
		ID: playerID,
	}

	// Logic to add the player to the session
	if session, exists := gs.Store.GetSession(sessionId); exists {
		// Implement the logic to add a player to this session.
		// This might involve creating a new Player instance and adding it to the session's Players map.
		session.AddPlayer(player)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully joined the session"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
	}
}
