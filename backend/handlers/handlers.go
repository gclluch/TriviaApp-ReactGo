package handlers

import (
	"github.com/gclluch/captrivia_multiplayer/game"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, gameServer *game.GameServer) {
	router.POST("/game/start", gameServer.StartGameHandler)
	router.POST("/game/join", gameServer.JoinGameHandler)
	router.POST("/questions", gameServer.QuestionsHandler)
	router.POST("/answer", gameServer.AnswerHandler)
	router.POST("/game/end", gameServer.EndGameHandler)
	router.POST("/player/finished", gameServer.MarkPlayerFinishedHandler)
	router.GET("/final-scores/:sessionId", gameServer.FinalScoresHandler)
	router.GET("/ws", gameServer.WebSocketEndpoint)
}
