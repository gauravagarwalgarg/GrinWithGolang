// Package main is the entry point for the clean architecture backend.
//
// LEARNING NOTES:
// - This file wires together ALL layers: bootstrap → route → controller → usecase → repository
// - The only place where concrete implementations are assembled (Composition Root)
// - After this point, every layer depends only on interfaces (Dependency Inversion Principle)
//
// For C++ devs: Think of this as your main() that constructs the object graph
// before passing control to the HTTP server event loop.
package main

import (
	"time"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/api/route"
	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Bootstrap: load config + connect to MongoDB
	app := bootstrap.App()
	env := app.Env
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	// 2. Derive timeout from config (seconds → time.Duration)
	timeout := time.Duration(env.ContextTimeout) * time.Second

	// 3. Create the Gin HTTP engine
	router := gin.Default()

	// 4. Wire all routes (public + protected with JWT middleware)
	route.Setup(env, timeout, db, router)

	// 5. Start the server
	router.Run(env.ServerAddress)
}
