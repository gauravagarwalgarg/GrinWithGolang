package domain

import "github.com/golang-jwt/jwt/v4"

// JwtCustomClaims holds the claims embedded in the access token.
// Name + ID are custom; StandardClaims provides exp, iat, etc.
type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	jwt.StandardClaims
}

// JwtCustomRefreshClaims holds minimal claims for refresh tokens.
// Only ID is needed refresh tokens are short-lived and just identify the user.
type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
