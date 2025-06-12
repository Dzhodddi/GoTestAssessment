package main

import (
	"FIDOtestBackendApp/internal/db"
	"FIDOtestBackendApp/internal/env"
	"FIDOtestBackendApp/internal/store"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

const version = "0.0.1"

//	@title			Golang engineer test assessment - the Spy Cat Agency
//	@description	API for Golang engineer test assessment - the Spy Cat Agency
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath		/v1
//
// @description	Golang engineer test assessment - the Spy Cat Agency
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("APP_ENV", "development"),
		db: dbConfig{
			addr:               env.GetString("DB_ADDR", "postgresql://postgres:adminadmin@localhost:5432/testGo?sslmode=disable"),
			maxOpenConnections: 10,
			maxIdleConnections: 10,
			maxIdleTime:        env.GetString("maxIdleTime", "15m"),
		},
	}

	// Logger init
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database init
	database, err := db.New(cfg.db.addr,
		cfg.db.maxOpenConnections,
		cfg.db.maxIdleConnections,
		cfg.db.maxIdleTime,
	)
	defer database.Close()
	if err != nil {
		logger.Fatal(err)
	}

	// Storage init
	storage := store.NewStorage(database)
	app := &application{
		config: cfg,
		logger: logger,
		store:  storage,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
