package repositories

import (
	"arfdev/chat/internal/domain/entities"
	"context"
)

type UserRepositoryInterface interface {
	FindByEmail(ctx context.Context, email string) (entities.Mst_users, error)
	Insert(ctx context.Context, payload entities.Mst_otp) error
	FindOTP(ctx context.Context, otpCode, signId string) (entities.Mst_otp, error)
}
