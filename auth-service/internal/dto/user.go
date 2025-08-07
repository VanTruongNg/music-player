package dto

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

type UserRegisterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
}

type UserLoginResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
	CreatedAt   string `json:"createdAt"`
	AccessToken string `json:"accessToken"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	CreatedAt string `json:"createdAt"`
}
