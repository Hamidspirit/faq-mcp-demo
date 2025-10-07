package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Hamidspirit/faq-mcp-demo.git/internal/models"
	"github.com/Hamidspirit/faq-mcp-demo.git/internal/mygenai"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	geminiClient *genai.GeminiClient
}

func NewChatHandler(geminiClient *genai.GeminiClient) *ChatHandler {
	return &ChatHandler{
		geminiClient: geminiClient,
	}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req models.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: "Message cannot be empty",
		})
		return
	}

	log.Printf("Received chat message: %s", req.Message)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get response from Gemini
	response, err := h.geminiClient.Chat(ctx, req.Message)
	if err != nil {
		log.Printf("Error getting chat response: %v", err)
		c.JSON(http.StatusInternalServerError, models.ChatResponse{
			Error: "Failed to generate response: " + err.Error(),
		})
		return
	}

	log.Printf("Generated response: %s", response)

	c.JSON(http.StatusOK, models.ChatResponse{
		Response: response,
	})
}

func (h *ChatHandler) HandleChatStream(c *gin.Context) {
	// TODO: Implement streaming response if needed
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Streaming not implemented yet",
	})
}
