package impl

import (
	"arfdev/chat/api/controllers"
	"arfdev/chat/internal/application/services"

	"github.com/gin-gonic/gin"
)

type authControllerImpl struct {
	authService services.AuthServiceInterface
}

func NewAuthControllerImpl(authService services.AuthServiceInterface) controllers.AuthControllerInterface {
	return &authControllerImpl{
		authService: authService,
	}
}

func (c *authControllerImpl) CheckPhone(ctx *gin.Context) {}

func (c *authControllerImpl) SendOtp(ctx *gin.Context) {}

func (c *authControllerImpl) ValidateOtp(ctx *gin.Context) {}

func (c *authControllerImpl) Register(ctx *gin.Context) {}

func (c *authControllerImpl) Login(ctx *gin.Context) {}

func (c *authControllerImpl) RefreshToken(ctx *gin.Context) {}

func (c *authControllerImpl) Logout(ctx *gin.Context) {}
