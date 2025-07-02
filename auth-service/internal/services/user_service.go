package services

import (
	"auth-service/internal/domain"
	"auth-service/internal/dto"
	repo "auth-service/internal/repositories"
	tokenmanager "auth-service/internal/services/TokenManager"
	"auth-service/internal/utils/jwt"
	"context"
	"time"
)

// UserService defines business logic for user resource.
type UserService interface {
	Register(ctx context.Context, req *dto.UserCreateRequest) (*domain.User, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*domain.User, string, string, error)
}

type userService struct {
	userRepo     repo.UserRepository
	jwtService   jwt.JWTService
	tokenManager tokenmanager.TokenManager
}

// NewUserService injects UserRepository, JWTService, and TokenManager for business logic.
func NewUserService(userRepo repo.UserRepository, jwtService jwt.JWTService, tokenManager tokenmanager.TokenManager) UserService {
	return &userService{
		userRepo:     userRepo,
		jwtService:   jwtService,
		tokenManager: tokenManager,
	}
}

func (s *userService) Register(ctx context.Context, req *dto.UserCreateRequest) (*domain.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrEmailExists
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) Login(ctx context.Context, req *dto.UserLoginRequest) (*domain.User, string, string, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", "", err
	}
	if existingUser == nil {
		return nil, "", "", domain.ErrInvalidCredentials
	}
	if !existingUser.CheckPassword(req.Password) {
		return nil, "", "", domain.ErrInvalidCredentials
	}

	now := time.Now().UTC().Format(time.RFC3339)
	existingUser.LastLoginAt = &now
	updatedUser, err := s.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(ctx, updatedUser.ID)
	if err != nil {
		return nil, "", "", err
	}

	return updatedUser, accessToken, refreshToken, nil
}
