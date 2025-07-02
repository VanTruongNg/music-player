package dto

// UserCreateRequest defines the expected payload for creating a user.
type UserCreateRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=64"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	FullName        string `json:"fullName" binding:"omitempty,max=64"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// UserRegisterResponse defines the response payload for a registered user.
type UserRegisterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
}

// UserLoginResponse defines the response payload for a successful login
// Includes user info and JWT tokens
// All fields are required for client authentication flows
type UserLoginResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"fullName"`
	CreatedAt    string `json:"createdAt"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
