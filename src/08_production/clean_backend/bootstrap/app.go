package bootstrap

import "github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/mongo"

// Application is the top-level struct that holds the app's dependencies.
// It's created once at startup and passed down to route setup.
//
// LEARNING NOTES:
// - This is the "Composition Root" where all dependencies are assembled
// - CloseDBConnection uses defer in main() for graceful shutdown
type Application struct {
	Env   *Env
	Mongo mongo.Client
}

// App initializes the application: loads config, connects to DB.
func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	return *app
}

// CloseDBConnection gracefully closes the MongoDB connection.
func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
