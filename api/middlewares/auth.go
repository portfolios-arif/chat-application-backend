package middlewares

import (
	"arfdev/chat/internal/application/dtos/responses"
	"arfdev/chat/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func abortHandler(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, responses.NewResponse(http.StatusUnauthorized, false, message, nil))
	ctx.Abort()
	return
}

func AuthMiddleware(jwtHelper *helpers.JWTHelper, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if "" == authHeader {
			abortHandler(ctx, "Unauthorized")
			return
		}

		token, err := jwtHelper.ExtractBearerToken(authHeader)
		if err != nil {
			abortHandler(ctx, err.Error())
			return
		}

		claims, err := jwtHelper.ValidateToken(token)
		if err != nil {
			abortHandler(ctx, err.Error())
			return
		}

		// Check redis
		find := rdb.Get(ctx.Request.Context(), "access-"+claims.UserID)
		if find.Err() != nil || find.Val() == "" {
			abortHandler(ctx, "Unauthorized")
			return
		}

		ctx.Set("UserID", claims.UserID)
		ctx.Next()
	}
}
