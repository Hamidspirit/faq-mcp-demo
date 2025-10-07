package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getTestRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")

	addTestRoute(v1)
}

func addTestRoute(rg *gin.RouterGroup) {
	test := rg.Group("test")

	test.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "this better be working")
	})
}
