package helpers

import (
	"arfdev/chat/internal/application/dtos/responses"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type serviceFn[T any] func(ctx context.Context, payload T) responses.APIBaseResponse

func POSTController[T any](service serviceFn[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload T

		if err := ctx.BindJSON(&payload); err != nil {
			response := responses.NewResponse(http.StatusBadRequest, false, "Invalid Body Request", nil)
			ctx.AbortWithStatusJSON(response.Code, response)
			return
		}

		validate := validator.New()
		if err := validate.Struct(&payload); err != nil {
			errValidation := err.(validator.ValidationErrors)
			response := responses.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil)
			ctx.AbortWithStatusJSON(response.Code, response)
			return
		}

		response := service(ctx.Request.Context(), payload)
		ctx.JSON(response.Code, response)
		return
	}
}
