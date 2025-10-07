package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "FAQ Chatbot API",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"home": "/",
			"test": "/test",
			"chat": "/chat",
		},
	})
}
