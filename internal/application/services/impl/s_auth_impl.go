package impl

import (
	"arfdev/chat/internal/application/dtos/requests"
	"arfdev/chat/internal/application/dtos/responses"
	"arfdev/chat/internal/application/services"
	"arfdev/chat/internal/domain/entities"
	"arfdev/chat/internal/domain/repositories"
	"arfdev/chat/pkg/helpers"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

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

func (s *authServiceImpl) CheckEmail(ctx context.Context, payload requests.CheckEmailRequestPayload) responses.APIBaseResponse {
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	user, err := s.repo.FindByEmail(ctx, payload.Email)
	if !payload.IsLogin {
		if user.Email == "" {
			return responses.NewResponse(http.StatusOK, true, "Email can be used for register", nil)
		}
		if err != nil {
			return responses.NewResponse(http.StatusInternalServerError, false, err.Error(), nil)
		}
		return responses.NewResponse(http.StatusBadRequest, false, "Email already exists", nil)
	}

	if user.Email != "" {
		response := responses.CheckPhoneIsLoginResponsePayload{
			DeviceID: user.DeviceID,
		}
		return responses.NewResponse(http.StatusOK, true, "User exists", response)
	}

	return responses.NewResponse(http.StatusBadRequest, false, "User does not exist", nil)
}

func (s *authServiceImpl) SendOtp(ctx context.Context, payload requests.OTPRequestPayload) responses.APIBaseResponse {
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	signatureID, err := helpers.GenerateUniqueID()
	otp, err := helpers.GenerateOTP()
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	data := &entities.Mst_otp{
		OTPCode:     otp,
		DeviceID:    payload.DeviceID,
		SignatureID: signatureID,
		ExpiredAt:   time.Now().Add(3 * time.Minute),
	}

	err = s.repo.Insert(ctx, *data)
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = helpers.SendEmailOTP(ctx, otp, payload.Email)
	if err != nil {
		switch {
		case ctx.Err() == context.DeadlineExceeded:
			return responses.NewResponse(http.StatusGatewayTimeout, false, "Request timeout: OTP email sending failed", nil)
		case strings.Contains(err.Error(), "authentication failed"):
			log.Printf("SMTP authentication error: %v", err)
			return responses.NewResponse(http.StatusInternalServerError, false, "Failed to send OTP email", nil)
		case strings.Contains(err.Error(), "connection failed"):
			log.Printf("SMTP connection error: %v", err)
			return responses.NewResponse(http.StatusServiceUnavailable, false, "Email service temporarily unavailable", nil)
		default:
			log.Printf("Email sending error: %v", err)
			return responses.NewResponse(http.StatusInternalServerError, false, "Failed to send OTP email", nil)
		}
	}

	response := responses.OTPResponsePayload{
		SignatureID: signatureID,
	}

	return responses.NewResponse(http.StatusOK, true, "Success request OTP", response)
}
