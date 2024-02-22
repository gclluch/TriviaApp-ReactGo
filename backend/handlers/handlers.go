package handlers

import (
	"github.com/gclluch/captrivia_multiplayer/game"
	"github.com/gin-gonic/gin"
)

// RegisterHandlers sets up the routing for the game server's API.
func RegisterHandlers(router *gin.Engine, gameServer *game.GameServer) {
	// Setup a group for game-related routes
	gameRoutes := router.Group("/game")
	{
		gameRoutes.POST("/start", gameServer.StartGameHandler)          // Start a new game session
		gameRoutes.POST("/join/:sessionId", gameServer.JoinGameHandler) // Join an existing game session
		gameRoutes.GET("/end/:sessionId", gameServer.EndGameHandler)    // End a game session
	}

	// Questions and answers handling
	router.GET("/questions/:sessionId", gameServer.QuestionsHandler) // Retrieve questions for the game
	router.POST("/answer", gameServer.AnswerHandler)                 // Submit an answer

	// Player status updates
	router.POST("/player/finished", gameServer.MarkPlayerFinishedHandler) // Mark a player as finished

	// Retrieve the final scores after a game session
	router.GET("/final-scores/:sessionId", gameServer.FinalScoresHandler)

	// WebSocket endpoint for real-time interactions
	router.GET("/ws", gameServer.WebSocketEndpoint)

	// Apply middleware for error handling (hypothetical example)
	router.Use(ErrorHandlingMiddleware())
}

// Centralized error handling.
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
		}
	}
}
