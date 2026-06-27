// Package middleware provides HTTP middleware for the Gin router.
//
// LEARNING NOTES:
// - Middleware in Gin is a function that returns gin.HandlerFunc
// - It runs BEFORE the actual handler and can short-circuit with c.Abort()
// - This JWT middleware:
//   1. Extracts the "Bearer <token>" from Authorization header
//   2. Validates the token signature and expiration
//   3. Extracts user ID from claims and sets it in Gin context
//   4. Downstream handlers read userID via c.GetString("x-user-id")
package middleware

import (
	"net/http"
	"strings"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/domain"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/internal/tokenutil"
	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware protects routes by requiring a valid access token.
func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")

		if len(t) == 2 {
			authToken := t[1]
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if authorized {
				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}
				// Store user ID in context for downstream handlers
				c.Set("x-user-id", userID)
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}
}
