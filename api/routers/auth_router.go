package routers

import (
	"arfdev/chat/api/controllers"
	"arfdev/chat/api/controllers/impl"
	"arfdev/chat/api/middlewares"
	"arfdev/chat/internal/application/services"
	"arfdev/chat/pkg/helpers"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type authRoute struct {
	authController controllers.AuthControllerInterface
}

func NewAuthRoute(authService services.AuthServiceInterface) *authRoute {
	return &authRoute{
		authController: impl.NewAuthControllerImpl(authService),
	}
}

func (r *authRoute) Setup(rg *gin.RouterGroup, rdb *redis.Client) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	accessExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRES_IN"))
	accessExpDuration := time.Duration(accessExp) * time.Hour
	jwtHelper := helpers.NewJWTHelper(secret, accessExpDuration)

	auth := rg.Group("/auth")
	{
		auth.POST("/check-email", r.authController.CheckEmail)
		auth.POST("/otp", r.authController.SendOtp)
		auth.POST("/validate-otp", r.authController.ValidateOtp)
		auth.POST("/register", r.authController.Register)
		auth.POST("/login", r.authController.Login)
		auth.POST("/refresh-token", middlewares.AuthMiddleware(jwtHelper, rdb), r.authController.RefreshToken)
		auth.GET("/logout", middlewares.AuthMiddleware(jwtHelper, rdb), r.authController.Logout)
	}
}
