package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"faq-chatbot/internal/models"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type FAQServer struct {
	server   *server.MCPServer
	faqData  *models.FAQDatabase
	dataPath string
}

func NewFAQServer(dataPath string) (*FAQServer, error) {
	fs := &FAQServer{
		dataPath: dataPath,
	}

	// Load FAQ data
	if err := fs.loadFAQData(); err != nil {
		return nil, fmt.Errorf("failed to load FAQ data: %w", err)
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"FAQ Database Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	fs.server = mcpServer

	// Register tools
	fs.registerTools()

	return fs, nil
}

func (fs *FAQServer) loadFAQData() error {
	data, err := os.ReadFile(fs.dataPath)
	if err != nil {
		return err
	}

	fs.faqData = &models.FAQDatabase{}
	return json.Unmarshal(data, fs.faqData)
}

func (fs *FAQServer) registerTools() {
	// Tool 1: Search FAQs by keyword
	fs.server.AddTool(mcp.Tool{
		Name:        "search_faqs",
		Description: "Search FAQs by keyword in questions, answers, or tags. Returns relevant FAQ entries.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"keyword": map[string]interface{}{
					"type":        "string",
					"description": "Keyword to search for in FAQs",
				},
			},
			Required: []string{"keyword"},
		},
	}, fs.handleSearchFAQs)

	// Tool 2: Get FAQ by category
	fs.server.AddTool(mcp.Tool{
		Name:        "get_faqs_by_category",
		Description: "Get all FAQs for a specific category (e.g., general, shipping, payment, account, products)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"category": map[string]interface{}{
					"type":        "string",
					"description": "Category name to filter FAQs",
				},
			},
			Required: []string{"category"},
		},
	}, fs.handleGetByCategory)

	// Tool 3: Get all categories
	fs.server.AddTool(mcp.Tool{
		Name:        "get_categories",
		Description: "Get list of all available FAQ categories",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, fs.handleGetCategories)

	// Tool 4: Get FAQ by ID
	fs.server.AddTool(mcp.Tool{
		Name:        "get_faq_by_id",
		Description: "Get a specific FAQ by its ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "FAQ ID",
				},
			},
			Required: []string{"id"},
		},
	}, fs.handleGetByID)
}

func (fs *FAQServer) handleSearchFAQs(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	keyword, ok := arguments["keyword"].(string)
	if !ok {
		return mcp.NewToolResultError("keyword must be a string"), nil
	}

	keyword = strings.ToLower(keyword)
	var results []models.FAQ

	for _, faq := range fs.faqData.FAQs {
		if strings.Contains(strings.ToLower(faq.Question), keyword) ||
			strings.Contains(strings.ToLower(faq.Answer), keyword) ||
			containsTag(faq.Tags, keyword) {
			results = append(results, faq)
		}
	}

	resultJSON, _ := json.MarshalIndent(results, "", "  ")
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (fs *FAQServer) handleGetByCategory(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	category, ok := arguments["category"].(string)
	if !ok {
		return mcp.NewToolResultError("category must be a string"), nil
	}

	category = strings.ToLower(category)
	var results []models.FAQ

	for _, faq := range fs.faqData.FAQs {
		if strings.ToLower(faq.Category) == category {
			results = append(results, faq)
		}
	}

	resultJSON, _ := json.MarshalIndent(results, "", "  ")
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (fs *FAQServer) handleGetCategories(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	categories := make(map[string]bool)
	for _, faq := range fs.faqData.FAQs {
		categories[faq.Category] = true
	}

	var categoryList []string
	for cat := range categories {
		categoryList = append(categoryList, cat)
	}

	resultJSON, _ := json.MarshalIndent(categoryList, "", "  ")
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (fs *FAQServer) handleGetByID(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	id, ok := arguments["id"].(string)
	if !ok {
		return mcp.NewToolResultError("id must be a string"), nil
	}

	for _, faq := range fs.faqData.FAQs {
		if faq.ID == id {
			resultJSON, _ := json.MarshalIndent(faq, "", "  ")
			return mcp.NewToolResultText(string(resultJSON)), nil
		}
	}

	return mcp.NewToolResultError(fmt.Sprintf("FAQ with ID %s not found", id)), nil
}

func (fs *FAQServer) GetServer() *server.MCPServer {
	return fs.server
}

func containsTag(tags []string, keyword string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), keyword) {
			return true
		}
	}
	return false
}
