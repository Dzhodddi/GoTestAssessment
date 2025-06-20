package main

import (
	"FIDOtestBackendApp/internal/db"
	"FIDOtestBackendApp/internal/env"
	"FIDOtestBackendApp/internal/graphql"
	"FIDOtestBackendApp/internal/store"
	"FIDOtestBackendApp/internal/store/cache"
	"errors"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

const version = "0.0.1"

var (
	ConflictError = errors.New("conflict")
)

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
		redisConfig: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       0,
			enabled:  true,
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

	//redis
	var cacheRedis *redis.Client
	if cfg.redisConfig.enabled {
		cacheRedis = cache.NewRedisClient(cfg.redisConfig.addr, cfg.redisConfig.password, cfg.redisConfig.db)
	}
	cacheStorage := cache.NewRedisStorage(cacheRedis)
	graphqlStorage := graphql.NewGPQLStorage()
	app := &application{
		config:         cfg,
		logger:         logger,
		store:          storage,
		cacheStorage:   cacheStorage,
		graphqlStorage: graphqlStorage,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
