package services

import (
	resDto "arfdev/chat/internal/application/dtos/requests"
	reqDto "arfdev/chat/internal/application/dtos/responses"
	"context"
)

type AuthServiceInterface interface {
	CheckPhone(ctx context.Context, payload resDto.CheckPhoneRequestPayload) reqDto.APIBaseResponse
}
