package services

import (
	"auth-service/internal/domain"
	"auth-service/internal/dto"
	repo "auth-service/internal/repositories"
	tokenmanager "auth-service/internal/services/TokenManager"
	"auth-service/internal/utils/jwt"
	"context"
	"log"
	"time"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetMe(ctx context.Context, userID string) (*domain.User, error)
	Register(ctx context.Context, req *dto.UserCreateRequest) (*domain.User, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*domain.User, string, string, error)
	RefreshToken(ctx context.Context, token string) (string, string, error)
	Logout(ctx context.Context, sid string) error
}

type userService struct {
	userRepo       repo.UserRepository
	jwtService     jwt.JWTService
	tokenManager   tokenmanager.TokenManager
	eventPublisher EventPublisher
}

func NewUserService(
	userRepo repo.UserRepository,
	jwtService jwt.JWTService,
	tokenManager tokenmanager.TokenManager,
	eventPublisher EventPublisher,
) UserService {
	return &userService{
		userRepo:       userRepo,
		jwtService:     jwtService,
		tokenManager:   tokenManager,
		eventPublisher: eventPublisher,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, domain.ErrNoUsersFound
	}
	return users, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
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

	// Publish user registered event (non-blocking, only log warning on failure)
	if err := s.eventPublisher.PublishUserRegistered(ctx, createdUser); err != nil {
		log.Printf("[WARN] Failed to publish user.registered event: %v", err)
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

	accessToken, refreshToken, err := s.tokenManager.IssueInitialTokens(ctx, updatedUser.ID)
	if err != nil {
		return nil, "", "", err
	}

	return updatedUser, accessToken, refreshToken, nil
}

func (s *userService) RefreshToken(ctx context.Context, token string) (string, string, error) {
	claims, err := s.jwtService.VerifyRefreshToken(token)
	if err != nil {
		return "", "", err
	}

	newAccessToken, newRefreshToken, err := s.tokenManager.RefreshToken(ctx, claims)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

func (s *userService) Logout(ctx context.Context, sid string) error {
	return s.tokenManager.RevokeSession(ctx, sid)
}
