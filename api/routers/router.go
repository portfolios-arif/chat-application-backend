package routers

import (
	"arfdev/chat/api/middlewares"
	"arfdev/chat/config"
	resDto "arfdev/chat/internal/application/dtos/responses"
	sImpl "arfdev/chat/internal/application/services/impl"
	"arfdev/chat/internal/infrastructures/persistence"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func welcomeHandler(ctx *gin.Context) {
	response := resDto.NewResponse(
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
	rdb := config.NewRedisClient()
	db := config.NewDB()

	// Middlewares
	router.Use(middlewares.NewSecurityHeaderMiddleware(nil))

	// Repositories
	userRepo := persistence.NewUserRepository(db)

	// Services
	authService := sImpl.NewAuthServiceImpl(userRepo, rdb)

	// Routers
	authRoute := NewAuthRoute(authService)

	v1 := router.Group(basePath)
	v1.Use(middlewares.TimeoutMiddleware(10 * time.Second))
	{
		v1.GET("/", welcomeHandler)
		authRoute.Setup(v1, rdb)
	}

	return router
}
