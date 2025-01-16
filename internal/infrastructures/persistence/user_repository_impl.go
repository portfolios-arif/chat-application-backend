package persistence

import (
	"arfdev/chat/internal/domain/entities"
	"arfdev/chat/internal/domain/repositories"
	"context"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepositoryInterface {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (rp *UserRepositoryImpl) FindByPhonenumber(ctx context.Context, phonenumber string) (entities.Mst_users, error) {
	return entities.Mst_users{}, nil
}
