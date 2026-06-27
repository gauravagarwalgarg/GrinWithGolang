package domain

import "context"

// RefreshTokenRequest is the expected payload for POST /refresh.
type RefreshTokenRequest struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
}

// RefreshTokenResponse returns new tokens after a successful refresh.
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenUsecase defines the business contract for token refresh.
type RefreshTokenUsecase interface {
	GetUserByID(c context.Context, id string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
	ExtractIDFromToken(requestToken string, secret string) (string, error)
}
