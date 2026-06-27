package domain

import "context"

// LoginRequest is the expected payload for POST /login.
// `binding:"required"` tells Gin to reject requests missing these fields.
type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginUsecase defines the business contract for login operations.
type LoginUsecase interface {
	GetUserByEmail(c context.Context, email string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
}
