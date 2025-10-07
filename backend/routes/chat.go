package routes

import (
	"context"
	"net/http"
	"os"

	"github.com/Hamidspirit/faq-mcp-demo.git/mygenai"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/genai"
)

var APIKEY = "apikey"

func getChatRoutes(r *gin.Engine) {
	chat := r.Group("v1")

	addChatRoutes(chat)
}

func addChatRoutes(rg *gin.RouterGroup) {
	chat := rg.Group("/chat")

	chat.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "This is chat route")
	})

	chat.POST("/chat", handleChat)
}

// Send and retrieve user request and LLM response
func handleChat(c *gin.Context) {
	var req mygenai.ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv(APIKEY),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Gemini client: " + err.Error()})
		return
	}

	defer client.Close()

	// This supposed to mimic MCP tool. i am using googles official sdk
	functionDecl := &genai.FunctionDeclaration{
		Name:        "search_faq",
		Description: "Search internal FAQs for user questionss",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"query": {
					Type:        genai.TypeString,
					Description: "The user's questions",
				},
			},
			Required: []string{"query"},
		},
	}

	// Configure model
	model := client.Models("gemini-2.5-flash")
}
