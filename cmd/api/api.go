package main

import (
	"FIDOtestBackendApp/docs"
	"FIDOtestBackendApp/internal/env"
	"FIDOtestBackendApp/internal/store"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config config
	logger *zap.SugaredLogger
	store  store.Storage
}

type dbConfig struct {
	addr               string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf(env.GetString("apiURL", "localhost"), env.GetString("ADDR", ":8080"))
	docs.SwaggerInfo.BasePath = "/v1"
	srv := echo.New()
	srv.Any("/*", echo.WrapHandler(mux))
	err := srv.Start(app.config.addr)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) mount() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			err := next(c)
			stop := time.Now()
			latency := stop.Sub(start)
			app.logger.Infof("%s %s %d %v", req.Method, req.RequestURI, res.Status, latency)
			return err
		}
	})
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	v1 := e.Group("/v1")
	v1.GET("/health", app.healthCheckHandler)
	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	cats := v1.Group("/spycat")
	cats.POST("", app.createCatHandler)
	cats.DELETE("/:id", app.deleteCatHandler)
	cats.GET("/:id", app.getCatByIDHandler)
	cats.PATCH("/:id", app.updateCatHandler)
	cats.GET("", app.getPaginatedCatListHandler)
	return e
}
