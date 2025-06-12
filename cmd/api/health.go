package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// HealthCheckAPI  godoc
//
//	@Summary		Health check
//	@Description	Health check
//	@Tags			health
//	@Produce		json
//	@Success		204	{object}	string	"OK"
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (app *application) healthCheckHandler(c echo.Context) error {
	data := map[string]string{
		"version": version,
		"status":  "OK",
		"env":     app.config.env,
	}
	return c.JSON(http.StatusOK, data)
}
