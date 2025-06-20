package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// List of Cats
//
//	@Summary		List of cats
//	@Description	List of cats
//	@Tags			spycat
//	@Success		200	{object}	string
//	@Failure		500	{object}	error
//	@Router			/ql [get]
func (app *application) getListOfCatsQL(c echo.Context) error {
	data := app.graphqlStorage.Cat.GetListOfCats()

	if len(data.Errors) > 0 {
		return c.JSON(http.StatusInternalServerError, data.Errors)
	}

	return c.JSON(http.StatusOK, data)
}
