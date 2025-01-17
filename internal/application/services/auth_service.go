package services

import (
	resDto "arfdev/chat/internal/application/dtos/requests"
	reqDto "arfdev/chat/internal/application/dtos/responses"
	"context"
)

type AuthServiceInterface interface {
	CheckEmail(ctx context.Context, payload resDto.CheckEmailRequestPayload) reqDto.APIBaseResponse
	SendOtp(ctx context.Context, payload resDto.OTPRequestPayload) reqDto.APIBaseResponse
}
