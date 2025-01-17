package middlewares

import (
	"arfdev/chat/internal/application/dtos/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan bool, 1)

		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			response := responses.NewResponse(http.StatusGatewayTimeout, false, "Request Timeout", nil)
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

	}
}
