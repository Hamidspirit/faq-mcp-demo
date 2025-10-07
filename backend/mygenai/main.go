package mygenai

import "strings"

type ChatRequest struct {
	Message string `json:"msg"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// Dummy data 07/10/2025
var faqs = map[string]string{
	"billing":  "Our billing is monthly via Stripe. Contact support for issues.",
	"setup":    "To set up, run `npm install` and `npm start`. See docs for details.",
	"features": "Key features: real-time chat, LLM integration, and custom tools.",
	"pricing":  "Free tier: 100 msgs/day. Pro: $10/mo unlimited.",
	"default":  "Sorry, I don't have info on that. Try asking about billing, setup, or pricing!",
}

// NOTE: Replace with the real deal
func searchFaq(query string) string {
	// Simple keyword search
	lowerQuery := strings.ToLower(query)
	for key := range faqs {
		if strings.Contains(lowerQuery, key) {
			return faqs[key]
		}
	}

	return faqs["def"]
}
