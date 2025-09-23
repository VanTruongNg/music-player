package dto

type UserRegisterRequest struct {
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