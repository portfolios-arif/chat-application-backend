package routers

import (
	"arfdev/chat/api/middlewares"
	"arfdev/chat/internal/domain"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func welcomeHandler(ctx *gin.Context) {
	response := domain.NewResponse(
		http.StatusOK,
		true,
		"Chat Application - Welcome",
		nil,
	)
	ctx.JSON(http.StatusOK, response)
}

func NewRoute() *gin.Engine {
	basePath := os.Getenv("APP_BASEPATH")

	router := gin.Default()

	// Middlewares
	router.Use(middlewares.NewSecurityHeaderMiddleware(nil))

	v1 := router.Group(basePath)
	{
		v1.GET("/", welcomeHandler)
	}

	return router
}
