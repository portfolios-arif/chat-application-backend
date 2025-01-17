package persistence

import (
	"arfdev/chat/internal/domain/entities"
	"arfdev/chat/internal/domain/repositories"
	"context"
	"errors"

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

func (rp *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (entities.Mst_users, error) {
	var find entities.Mst_users
	result := rp.db.
		WithContext(ctx).
		Model(&entities.Mst_users{}).
		Where("email = ?", email).
		Find(&find)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			find.Email = ""
			return find, nil
		} else {
			return entities.Mst_users{}, result.Error
		}
	}
	return find, nil
}

func (rp *UserRepositoryImpl) Insert(ctx context.Context, payload entities.Mst_otp) error {
	result := rp.db.
		WithContext(ctx).
		Model(&entities.Mst_otp{}).
		Create(&payload)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
