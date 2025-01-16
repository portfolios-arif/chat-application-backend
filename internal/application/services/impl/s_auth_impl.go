package impl

import (
	"arfdev/chat/internal/application/dtos/requests"
	"arfdev/chat/internal/application/dtos/responses"
	"arfdev/chat/internal/application/services"
	"arfdev/chat/internal/domain/repositories"
	"context"

	"github.com/redis/go-redis/v9"
)

type authServiceImpl struct {
	repo repositories.UserRepositoryInterface
	rdb  *redis.Client
}

func NewAuthServiceImpl(repo repositories.UserRepositoryInterface, rdb *redis.Client) services.AuthServiceInterface {
	return &authServiceImpl{
		repo: repo,
		rdb:  rdb,
	}
}

func (s *authServiceImpl) CheckPhone(ctx context.Context, payload requests.CheckPhoneRequestPayload) responses.APIBaseResponse {
	return responses.APIBaseResponse{}
}
