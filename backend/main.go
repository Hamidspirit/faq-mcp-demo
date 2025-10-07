package main

import (
	"context"
	"log"

	"github.com/Hamidspirit/faq-mcp-demo.git/config"
	"github.com/Hamidspirit/faq-mcp-demo.git/internal/handlers"
	"github.com/Hamidspirit/faq-mcp-demo.git/internal/mcp"
	"github.com/Hamidspirit/faq-mcp-demo.git/internal/mygenai"
	"github.com/Hamidspirit/faq-mcp-demo.git/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	// Create MCP FAQ Server
	mcpServer, err := mcp.NewFAQServer(cfg.FAQDataPath)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}
	log.Println("MCP FAQ Server initialized")

	// Create Gemini client with MCP support
	ctx := context.Background()
	geminiClient, err := genai.NewGeminiClient(ctx, cfg.GeminiAPIKey, mcpServer)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer geminiClient.Close()
	log.Println("Gemini client initialized")

	// Create handlers
	chatHandler := handlers.NewChatHandler(geminiClient)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	// Setup routes
	routes.SetupRoutes(router, chatHandler)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
