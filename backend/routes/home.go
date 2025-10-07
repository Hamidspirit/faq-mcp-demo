package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHomeRoutes(r *gin.Engine) {
	home := r.Group("v1")

	addHomeroutes(home)
}

func addHomeroutes(rg *gin.RouterGroup) {
	home := rg.Group("/")

	home.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "This is home route")
	})
}
