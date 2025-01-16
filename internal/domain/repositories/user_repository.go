package repositories

import (
	"arfdev/chat/internal/domain/entities"
	"context"
)

type UserRepositoryInterface interface {
	FindByPhonenumber(ctx context.Context, phonenumber string) (entities.Mst_users, error)
}
