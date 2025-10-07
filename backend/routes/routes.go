package routes

import (
	"github.com/Hamidspirit/faq-mcp-demo.git/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {
	// Home route group
	homeGroup := router.Group("/")
	{
		homeGroup.GET("", handlers.HomeHandler)
	}

	// Test route group
	testGroup := router.Group("/test")
	{
		testGroup.GET("", handlers.TestHandler)
	}

	// Chat route group
	chatGroup := router.Group("/chat")
	{
		chatGroup.POST("", chatHandler.HandleChat)
		chatGroup.POST("/stream", chatHandler.HandleChatStream)
	}
}
