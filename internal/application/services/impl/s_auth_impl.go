package impl

import (
	"arfdev/chat/internal/application/dtos/requests"
	"arfdev/chat/internal/application/dtos/responses"
	"arfdev/chat/internal/application/services"
	"arfdev/chat/internal/domain/repositories"
	"context"
	"net/http"

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
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	user, err := s.repo.FindByPhonenumber(ctx, payload.Phonenumber)
	if !payload.IsLogin {
		if user.Phonenumber == "" {
			return responses.NewResponse(http.StatusOK, true, "Phone number can be used for register", nil)
		}
		if err != nil {
			return responses.NewResponse(http.StatusInternalServerError, false, err.Error(), nil)
		}
		return responses.NewResponse(http.StatusBadRequest, false, "Phone number already exists", nil)
	}

	if user.Phonenumber != "" {
		response := responses.CheckPhoneIsLoginResponsePayload{
			DeviceID: user.DeviceID,
		}
		return responses.NewResponse(http.StatusOK, true, "User exists", response)
	}

	return responses.NewResponse(http.StatusBadRequest, false, "User does not exist", nil)
}
