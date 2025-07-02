package services

import (
	"auth-service/internal/domain"
	"auth-service/internal/repositories"
	redisutil "auth-service/internal/utils/redis"
	"auth-service/internal/utils/twofa"
	"context"
	"errors"
	"time"
)

type TwoFAService struct {
	twoFAUtil *twofa.TwoFAUtil
	userRepo  repositories.UserRepository
	redisUtil *redisutil.RedisUtil
}

func NewTwoFAService(userRepo repositories.UserRepository, twoFAUtil *twofa.TwoFAUtil, redisUtil *redisutil.RedisUtil) *TwoFAService {
	return &TwoFAService{
		twoFAUtil: twoFAUtil,
		userRepo:  userRepo,
		redisUtil: redisUtil,
	}
}

func (s *TwoFAService) Setup2FA(ctx context.Context, userID string) (*twofa.SetupResult, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	if user.TwoFAEnabled {
		return nil, domain.ErrTwoFAEnabled
	}

	existingKey := "2fa:setup:" + userID
	var cached twofa.SetupResult
	if err := s.redisUtil.GetJSON(ctx, existingKey, &cached); err == nil && cached.Secret != "" {
		return &cached, nil
	}

	setupResult, err := s.twoFAUtil.GenerateSecret(user.Email)
	if err != nil {
		return nil, err
	}
	redisKey := "2fa:setup:" + userID
	err = s.redisUtil.SetJSON(ctx, redisKey, setupResult, 300*time.Second)
	if err != nil {
		return nil, err
	}
	return setupResult, nil
}

func (s *TwoFAService) Enable2FA(ctx context.Context, userID, code string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}
	redisKey := "2fa:setup:" + userID
	var setup twofa.SetupResult
	if err := s.redisUtil.GetJSON(ctx, redisKey, &setup); err != nil || setup.Secret == "" {
		return errors.New("2FA secret not found or expired")
	}
	if err := s.twoFAUtil.VerifyCode(setup.Secret, code); err != nil {
		return err
	}
	user.TwoFAEnabled = true
	user.TwoFASecret = setup.Secret
	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	_ = s.redisUtil.Delete(ctx, redisKey)
	return nil
}

func (s *TwoFAService) Verify2FA(ctx context.Context, userID, code string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}
	if !user.TwoFAEnabled || user.TwoFASecret == "" {
		return domain.ErrTwoFANotAvailable
	}
	return s.twoFAUtil.VerifyCode(user.TwoFASecret, code)
}

func (s *TwoFAService) Disable2FA(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return domain.ErrUserNotFound
	}
	user.TwoFAEnabled = false
	user.TwoFASecret = ""
	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}
