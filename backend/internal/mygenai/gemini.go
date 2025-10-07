package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Hamidspirit/faq-mcp-demo.git/internal/mcp"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client    *genai.Client
	model     *genai.GenerativeModel
	mcpServer *mcp.FAQServer
}

func NewGeminiClient(ctx context.Context, apiKey string, mcpServer *mcp.FAQServer) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	// Configure model
	model.SetTemperature(0.7)
	model.SetTopP(0.95)
	model.SetTopK(40)
	model.SetMaxOutputTokens(2048)

	// Set system instruction
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(`You are a helpful FAQ assistant. Your job is to answer user questions using the FAQ database tools available to you.

IMPORTANT INSTRUCTIONS:
1. Always use the available tools to search for relevant FAQs before answering
2. When a user asks a question, use search_faqs tool with relevant keywords
3. Base your answers on the FAQ data retrieved from the tools
4. If you find relevant FAQs, summarize them in a friendly, conversational way
5. If no relevant FAQs are found, politely say you don't have that information
6. Be concise but helpful
7. Use proper formatting for readability

Available tools:
- search_faqs: Search FAQs by keyword
- get_faqs_by_category: Get FAQs by category
- get_categories: List all categories
- get_faq_by_id: Get specific FAQ by ID`),
		},
	}

	// Configure tools/function calling
	tools := convertMCPToolsToGemini(mcpServer)
	model.Tools = tools

	gc := &GeminiClient{
		client:    client,
		model:     model,
		mcpServer: mcpServer,
	}

	return gc, nil
}

func (gc *GeminiClient) Chat(ctx context.Context, message string) (string, error) {
	session := gc.model.StartChat()
	session.History = []*genai.Content{}

	// Send user message
	resp, err := session.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// Handle function calls
	for {
		// Check if model wants to call functions
		part := resp.Candidates[0].Content.Parts[0]

		functionCall, ok := part.(genai.FunctionCall)
		if !ok {
			// No function call, return the text response
			if textPart, ok := part.(genai.Text); ok {
				return string(textPart), nil
			}
			return "I apologize, but I couldn't generate a proper response.", nil
		}

		// Execute the function call via MCP
		log.Printf("Function call: %s with args: %v", functionCall.Name, functionCall.Args)

		functionResult, err := gc.executeMCPTool(functionCall.Name, functionCall.Args)
		if err != nil {
			return "", fmt.Errorf("failed to execute tool: %w", err)
		}

		log.Printf("Function result: %s", functionResult)

		// Send function result back to model
		resp, err = session.SendMessage(ctx, genai.FunctionResponse{
			Name: functionCall.Name,
			Response: map[string]interface{}{
				"result": functionResult,
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to send function response: %w", err)
		}
	}
}

func (gc *GeminiClient) executeMCPTool(toolName string, args map[string]interface{}) (string, error) {
	server := gc.mcpServer.GetServer()

	// Convert args to JSON for MCP
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", err
	}

	var argsMap map[string]interface{}
	if err := json.Unmarshal(argsJSON, &argsMap); err != nil {
		return "", err
	}

	// Call the tool
	result, err := server.CallTool(toolName, argsMap)
	if err != nil {
		return "", err
	}

	// Extract text from result
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(map[string]interface{}); ok {
			if text, ok := textContent["text"].(string); ok {
				return text, nil
			}
		}
	}

	return "", fmt.Errorf("unexpected result format")
}

func convertMCPToolsToGemini(mcpServer *mcp.FAQServer) []*genai.Tool {
	server := mcpServer.GetServer()
	tools := server.ListTools()

	geminiTools := make([]*genai.Tool, 0)

	functionDeclarations := make([]*genai.FunctionDeclaration, 0)

	for _, tool := range tools.Tools {
		// Convert MCP tool schema to Gemini function declaration
		params := &genai.Schema{
			Type:       genai.TypeObject,
			Properties: make(map[string]*genai.Schema),
			Required:   tool.InputSchema.Required,
		}

		// Convert properties
		for propName, propValue := range tool.InputSchema.Properties {
			propMap, ok := propValue.(map[string]interface{})
			if !ok {
				continue
			}

			propType, _ := propMap["type"].(string)
			propDesc, _ := propMap["description"].(string)

			params.Properties[propName] = &genai.Schema{
				Type:        getGeminiType(propType),
				Description: propDesc,
			}
		}

		functionDeclarations = append(functionDeclarations, &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  params,
		})
	}

	geminiTools = append(geminiTools, &genai.Tool{
		FunctionDeclarations: functionDeclarations,
	})

	return geminiTools
}

func getGeminiType(mcpType string) genai.Type {
	switch mcpType {
	case "string":
		return genai.TypeString
	case "number":
		return genai.TypeNumber
	case "boolean":
		return genai.TypeBoolean
	case "array":
		return genai.TypeArray
	case "object":
		return genai.TypeObject
	default:
		return genai.TypeString
	}
}

func (gc *GeminiClient) Close() error {
	return gc.client.Close()
}
