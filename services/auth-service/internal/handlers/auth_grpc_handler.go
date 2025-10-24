package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/dto"
	"auth-service/internal/services"
	tokenmanager "auth-service/internal/services/TokenManager"
	"context"
	authv1 "music-player/api/proto/auth/v1"
	"time"

	"google.golang.org/grpc/metadata"
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

	clientIP := "unknown"
	userAgent := "unknown"

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ips := md.Get("x-client-ip"); len(ips) > 0 {
			clientIP = ips[0]
		}
		if uas := md.Get("x-user-agent"); len(uas) > 0 {
			userAgent = uas[0]
		}
	}

	newCtx := context.WithValue(ctx, tokenmanager.CtxKeyIP, clientIP)
	newCtx = context.WithValue(newCtx, tokenmanager.CtxKeyUserAgent, userAgent)

	loginReq := &dto.UserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, accessToken, refreshToken, err := h.userService.Login(newCtx, loginReq)
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

// GetUserProfile gets user profile information
func (h *AuthGRPCHandler) GetUserProfile(ctx context.Context, req *authv1.GetUserProfileRequest) (*authv1.GetUserProfileResponse, error) {

	if req.UserId == "" {
		return &authv1.GetUserProfileResponse{
			Success: false,
			Message: "User ID is required",
		}, nil
	}

	user, err := h.userService.GetMe(ctx, req.UserId)
	if err != nil {
		return &authv1.GetUserProfileResponse{
			Success: false,
			Message: "Failed to retrieve user profile",
		}, nil
	}

	return &authv1.GetUserProfileResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		User: &authv1.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
		},
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
	clientIP := "unknown"
	userAgent := "unknown"
	refreshToken := ""

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ips := md.Get("x-client-ip"); len(ips) > 0 {
			clientIP = ips[0]
		}
		if uas := md.Get("x-user-agent"); len(uas) > 0 {
			userAgent = uas[0]
		}
		if tokens := md.Get("refresh_token"); len(tokens) > 0 {
			refreshToken = tokens[0]
		}
	}

	newCtx := context.WithValue(ctx, tokenmanager.CtxKeyIP, clientIP)
	newCtx = context.WithValue(newCtx, tokenmanager.CtxKeyUserAgent, userAgent)

	if refreshToken == "" {
		return &authv1.RefreshTokenResponse{
			Success: false,
			Message: "Refresh token is required",
		}, nil
	}

	accessToken, newRefreshToken, err := h.userService.RefreshToken(newCtx, refreshToken)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.RefreshTokenResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.RefreshTokenResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	return &authv1.RefreshTokenResponse{
		Success:      true,
		Message:      "Token refreshed successfully",
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
	}, nil
}

// Logout handles user logout
func (h *AuthGRPCHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if req.Sid == "" {
		return &authv1.LogoutResponse{
			Success: false,
			Message: "Session ID is required",
		}, nil
	}

	err := h.userService.Logout(ctx, req.Sid)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.LogoutResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.LogoutResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

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

func (h *AuthGRPCHandler) SetupTwoFA(ctx context.Context, req *authv1.SetupTwoFARequest) (*authv1.SetupTwoFAResponse, error) {
	if req.UserId == "" {
		return &authv1.SetupTwoFAResponse{
			Success: false,
			Message: "User ID is required",
		}, nil
	}
	result, err := h.twoFAService.Setup2FA(ctx, req.UserId)
	if err != nil {
		return &authv1.SetupTwoFAResponse{
			Success: false,
			Message: "Failed to setup 2FA: " + err.Error(),
		}, nil
	}

	return &authv1.SetupTwoFAResponse{
		Success: true,
		Message: "2FA setup successful",
		Secret:  result.Secret,
		OtpUrl:  result.OTPURL,
	}, nil
}

// EnableTwoFA enables two-factor authentication
func (h *AuthGRPCHandler) EnableTwoFA(ctx context.Context, req *authv1.EnableTwoFARequest) (*authv1.EnableTwoFAResponse, error) {
	if req.UserId == "" || req.Code == "" {
		return &authv1.EnableTwoFAResponse{
			Success: false,
			Message: "User ID and code are required",
		}, nil
	}

	err := h.twoFAService.Enable2FA(ctx, req.UserId, req.Code)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.EnableTwoFAResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.EnableTwoFAResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	return &authv1.EnableTwoFAResponse{
		Success: true,
		Message: "2FA enabled successfully",
	}, nil
}

// DisableTwoFA disables two-factor authentication
func (h *AuthGRPCHandler) DisableTwoFA(ctx context.Context, req *authv1.DisableTwoFARequest) (*authv1.DisableTwoFAResponse, error) {
	if req.UserId == "" || req.Code == "" {
		return &authv1.DisableTwoFAResponse{
			Success: false,
			Message: "User ID and code are required",
		}, nil
	}

	err := h.twoFAService.Disable2FA(ctx, req.UserId, req.Code)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.DisableTwoFAResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.DisableTwoFAResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}
	return &authv1.DisableTwoFAResponse{
		Success: true,
		Message: "2FA disabled successfully",
	}, nil
}

// VerifyTwoFA verifies a 2FA code
func (h *AuthGRPCHandler) VerifyTwoFA(ctx context.Context, req *authv1.VerifyTwoFARequest) (*authv1.VerifyTwoFAResponse, error) {
	if req.UserId == "" || req.Code == "" {
		return &authv1.VerifyTwoFAResponse{
			Success: false,
			Message: "User ID and code are required",
		}, nil
	}
	err := h.twoFAService.Verify2FA(ctx, req.UserId, req.Code)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			return &authv1.VerifyTwoFAResponse{
				Success: false,
				Message: derr.Message,
			}, nil
		}
		return &authv1.VerifyTwoFAResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}
	return &authv1.VerifyTwoFAResponse{
		Success: true,
		Message: "2FA code is valid",
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
