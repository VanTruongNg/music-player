package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/dto"
	"auth-service/internal/services"
	"context"
	authv1 "music-player/api/proto/auth/v1"
	"time"
)

type AuthGRPCHandler struct {
	authv1.UnimplementedAuthServiceServer
	userService  services.UserService
	twoFAService services.TwoFAService
}

func NewAuthGRPCHandler(
	userService services.UserService,
	twoFAService services.TwoFAService,
) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		userService:  userService,
		twoFAService: twoFAService,
	}
}

func (h *AuthGRPCHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &authv1.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	loginReq := &dto.UserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call existing user service
	user, accessToken, refreshToken, err := h.userService.Login(ctx, loginReq)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.LoginResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.LoginResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	if req.TwoFaCode != "" && user.TwoFAEnabled {
		if err := h.twoFAService.Verify2FA(ctx, user.ID, req.TwoFaCode); err != nil {
			return &authv1.LoginResponse{
				Success: false,
				Message: "Invalid 2FA code",
			}, nil
		}
	} else if user.TwoFAEnabled {
		return &authv1.LoginResponse{
			Success: false,
			Message: "2FA code required",
		}, nil
	}

	return &authv1.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		User: &authv1.User{
			Id:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			FullName:     user.FullName,
			TwoFaEnabled: user.TwoFAEnabled,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *AuthGRPCHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &authv1.RegisterResponse{
			Success: false,
			Message: "Username, email, and password are required",
		}, nil
	}

	registerReq := &dto.UserCreateRequest{
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.Password,
		FullName:        req.FullName,
	}

	createdUser, err := h.userService.Register(ctx, registerReq)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.RegisterResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.RegisterResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	return &authv1.RegisterResponse{
		Success: true,
		Message: "Registration successful",
		UserId:  createdUser.ID,
	}, nil
}

func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	if req.AccessToken == "" {
		return &authv1.ValidateTokenResponse{
			Valid:   false,
			Message: "Access token is required",
		}, nil
	}

	return &authv1.ValidateTokenResponse{
		Valid:   true,
		Message: "Token is valid",
		User: &authv1.User{
			Id:       "placeholder_user_id",
			Username: "username",
			Email:    "user@example.com",
			FullName: "User Name",
		},
		ExpiresAt: 1234567890,
	}, nil
}

// RefreshToken handles token refresh
func (h *AuthGRPCHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	// TODO: Implement token refresh using tokenManager

	if req.RefreshToken == "" {
		return &authv1.RefreshTokenResponse{
			Success: false,
			Message: "Refresh token is required",
		}, nil
	}

	return &authv1.RefreshTokenResponse{
		Success:      true,
		Message:      "Token refreshed successfully",
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		ExpiresIn:    3600,
	}, nil
}

// Logout handles user logout
func (h *AuthGRPCHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	// TODO: Implement logout logic (revoke tokens)

	return &authv1.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}

// RevokeToken revokes a token
func (h *AuthGRPCHandler) RevokeToken(ctx context.Context, req *authv1.RevokeTokenRequest) (*authv1.RevokeTokenResponse, error) {
	// TODO: Implement token revocation

	return &authv1.RevokeTokenResponse{
		Success: true,
		Message: "Token revoked successfully",
	}, nil
}

// EnableTwoFA enables two-factor authentication
func (h *AuthGRPCHandler) EnableTwoFA(ctx context.Context, req *authv1.EnableTwoFARequest) (*authv1.EnableTwoFAResponse, error) {
	// TODO: Implement 2FA enable using twoFAService

	return &authv1.EnableTwoFAResponse{
		Success:   true,
		Message:   "2FA enabled successfully",
		QrCodeUrl: "placeholder_qr_url",
		SecretKey: "placeholder_secret",
	}, nil
}

// DisableTwoFA disables two-factor authentication
func (h *AuthGRPCHandler) DisableTwoFA(ctx context.Context, req *authv1.DisableTwoFARequest) (*authv1.DisableTwoFAResponse, error) {
	// TODO: Implement 2FA disable using twoFAService

	return &authv1.DisableTwoFAResponse{
		Success: true,
		Message: "2FA disabled successfully",
	}, nil
}

// VerifyTwoFA verifies a 2FA code
func (h *AuthGRPCHandler) VerifyTwoFA(ctx context.Context, req *authv1.VerifyTwoFARequest) (*authv1.VerifyTwoFAResponse, error) {
	// TODO: Implement 2FA verification using twoFAService

	return &authv1.VerifyTwoFAResponse{
		Valid:   true,
		Message: "2FA code is valid",
	}, nil
}

// GetUserProfile gets user profile information
func (h *AuthGRPCHandler) GetUserProfile(ctx context.Context, req *authv1.GetUserProfileRequest) (*authv1.GetUserProfileResponse, error) {
	// TODO: Implement get user profile using userService

	return &authv1.GetUserProfileResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		User: &authv1.User{
			Id:       req.UserId,
			Username: "username",
			Email:    "user@example.com",
			FullName: "User Name",
		},
	}, nil
}

// UpdateUserProfile updates user profile information
func (h *AuthGRPCHandler) UpdateUserProfile(ctx context.Context, req *authv1.UpdateUserProfileRequest) (*authv1.UpdateUserProfileResponse, error) {
	// TODO: Implement update user profile using userService

	return &authv1.UpdateUserProfileResponse{
		Success: true,
		Message: "User profile updated successfully",
		User: &authv1.User{
			Id:       req.UserId,
			Username: "username",
			Email:    req.Email,
			FullName: req.FullName,
		},
	}, nil
}
