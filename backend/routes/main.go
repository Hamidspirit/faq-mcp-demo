package routes

import "github.com/gin-gonic/gin"

func CreateRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func Run() {
	r := CreateRouter()
	getTestRoutes(r)

	getHomeRoutes(r)
	getChatRoutes(r)

	r.Run(":5000")
}
