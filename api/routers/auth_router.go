package routers

import (
	"arfdev/chat/api/controllers"
	"arfdev/chat/api/controllers/impl"
	"arfdev/chat/internal/application/services"

	"github.com/gin-gonic/gin"
)

type authRoute struct {
	authController controllers.AuthControllerInterface
}

func NewAuthRoute(authService services.AuthServiceInterface) *authRoute {
	return &authRoute{
		authController: impl.NewAuthControllerImpl(authService),
	}
}

func (r *authRoute) Setup(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/check-phonenumber", r.authController.CheckPhone)
		auth.POST("/otp", r.authController.SendOtp)
		auth.POST("/validate-otp", r.authController.ValidateOtp)
		auth.POST("/register", r.authController.Register)
		auth.POST("/login", r.authController.Login)
		auth.GET("/refresh-token", r.authController.RefreshToken)
		auth.GET("/logout", r.authController.Logout)
	}
}
