package repositories

import (
	"auth-service/internal/domain"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) findOneByField(ctx context.Context, field string, value any) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where(field+" = ?", value).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return nil, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return r.findOneByField(ctx, "id", id)
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.findOneByField(ctx, "username", username)
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findOneByField(ctx, "email", email)
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
