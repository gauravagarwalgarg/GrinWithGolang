// Package route wires controllers to HTTP paths and applies middleware.
//
// LEARNING NOTES:
// - Two router groups: public (no auth) and protected (JWT required)
// - Each route file creates its own repository → usecase → controller chain
// - This is where Dependency Injection happens manually (no framework needed in Go!)
//
// For C++ devs: This is like registering callbacks with an HTTP server,
// but with middleware acting as interceptors (similar to aspect-oriented programming).
package route

import (
	"time"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/api/middleware"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/bootstrap"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/mongo"
	"github.com/gin-gonic/gin"
)

// Setup registers all application routes on the Gin engine.
func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	// --- Public routes (no authentication required) ---
	publicRouter := gin.Group("")
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)

	// --- Protected routes (JWT middleware validates access token) ---
	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	NewProfileRouter(env, timeout, db, protectedRouter)
	NewTaskRouter(env, timeout, db, protectedRouter)
}
