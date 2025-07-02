package domain

import "net/http"

type DomainError struct {
	Code    string
	Message string
	Status  int
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrEmailExists    = &DomainError{"EMAIL_EXISTS", "email already registered", http.StatusConflict}
	ErrUsernameExists = &DomainError{"USERNAME_EXISTS", "username already registered", http.StatusConflict}
	ErrUserNotFound   = &DomainError{"USER_NOT_FOUND", "user not found", http.StatusNotFound}
	ErrTwoFAEnabled  = &DomainError{"TWO_FA_ENABLED", "2FA is already enabled for this user", http.StatusBadRequest}
	ErrTwoFANotAvailable = &DomainError{"TWO_FA_NOT_AVAILABLE", "2FA is not available for this user", http.StatusBadRequest}
	ErrInvalidCredentials = &DomainError{"INVALID_CREDENTIALS", "invalid email or password", http.StatusUnauthorized}
)
