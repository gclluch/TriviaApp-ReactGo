package handlers

import (
	"github.com/gclluch/captrivia_multiplayer/game"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, gameServer *game.GameServer) {
	router.POST("/game/start", gameServer.StartGameHandler)
	router.GET("/questions", gameServer.QuestionsHandler)
	router.POST("/answer", gameServer.AnswerHandler)
	router.POST("/game/end", gameServer.EndGameHandler)
	router.GET("/ws", gameServer.WebSocketEndpoint)
}
