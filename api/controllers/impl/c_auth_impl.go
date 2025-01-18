package impl

import (
	"arfdev/chat/api/controllers"
	"arfdev/chat/internal/application/dtos/requests"
	"arfdev/chat/internal/application/dtos/responses"
	"arfdev/chat/internal/application/services"
	"arfdev/chat/pkg/helpers"
	"net/http"

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

func (c *authControllerImpl) CheckEmail(ctx *gin.Context) {
	helpers.POSTController(c.authService.CheckEmail)(ctx)
}

func (c *authControllerImpl) SendOtp(ctx *gin.Context) {
	helpers.POSTController(c.authService.SendOtp)(ctx)
}

func (c *authControllerImpl) ValidateOtp(ctx *gin.Context) {
	helpers.POSTController(c.authService.VerifyOTP)(ctx)
}

func (c *authControllerImpl) Register(ctx *gin.Context) {
	var payload requests.RegisterRequestPayload

	if err := ctx.ShouldBind(&payload); err != nil {
		response := responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
		ctx.AbortWithStatusJSON(response.Code, response)
		return
	}

	res := c.authService.Register(ctx.Request.Context(), payload)
	ctx.JSON(res.Code, res)
}

func (c *authControllerImpl) Login(ctx *gin.Context) {
	helpers.POSTController(c.authService.Login)(ctx)
}

func (c *authControllerImpl) RefreshToken(ctx *gin.Context) {
	helpers.POSTController(c.authService.RefreshToken)(ctx)
}

func (c *authControllerImpl) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("UserID")
	if !exists {
		response := responses.NewResponse(http.StatusUnauthorized, false, "Unauthorized", nil)
		ctx.AbortWithStatusJSON(response.Code, response)
	}

	res := c.authService.Logout(ctx.Request.Context(), userID.(string))
	ctx.JSON(res.Code, res)
}
