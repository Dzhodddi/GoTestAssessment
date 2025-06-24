package main

import (
	"WorkAssigment/internal/store"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

var (
	ValidationError = errors.New("validation error")
)

type CreateCatPayload struct {
	Name       string `json:"name" validate:"required,max=200"`
	Experience int    `json:"year_of_experience" validate:"required,gte=1"`
	Breed      string `json:"breed" validate:"required,max=200,breed-exits"`
	Salary     int    `json:"salary" validate:"required,gte=1"`
}

type UpdateCatInfoPayload struct {
	Salary int `json:"salary" validate:"required,gte=0"`
}

// Create SpyCat
//
//	@Summary		Create spy cat
//	@Description	Create new spy cat
//	@Tags			spycat
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCatPayload	true	"SpyCat payload"
//	@Success		201		{object}	store.Cat
//	@Failure		400		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/spycat [post]
func (app *application) createCatHandler(c echo.Context) error {
	var payload CreateCatPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}
	breed, err := app.cacheStorage.Cats.Get(c.Request().Context(), payload.Breed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if !breed {
		if err = Validate.Struct(payload); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
		}
		err = app.cacheStorage.Cats.Set(c.Request().Context(), payload.Breed)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	spyCat := &store.Cat{
		Name:       payload.Name,
		Breed:      payload.Breed,
		Experience: payload.Experience,
		Salary:     payload.Salary,
	}

	if err := app.store.Cat.CreateSpyCat(c.Request().Context(), spyCat); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, spyCat)
}

// Delete cat godoc
//
//	@Summary		Delete spy cat
//	@Description	Delete spy cat by ID
//	@Tags			spycat
//	@Produce		json
//	@Param			id	path		int	true	"Cat ID"
//	@Success		204	{object}	nil
//	@Failure		422	{object}	error
//	@Failure		500	{object}	error
//	@Router			/spycat/{id} [delete]
func (app *application) deleteCatHandler(c echo.Context) error {
	cat, err := app.getCatByID(c)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, ValidationError):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	err = app.store.Cat.DeleteSpyCat(c.Request().Context(), cat.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusNoContent, nil)
}

// Get cat info godoc
//
//	@Summary		Get cat info
//	@Description	Get cat info by ID
//	@Tags			spycat
//	@Produce		json
//	@Param			id	path		int	true	"Cat ID"
//	@Success		200	{object}	store.Cat
//	@Failure		422	{object}	error
//	@Failure		500	{object}	error
//	@Router			/spycat/{id} [get]
func (app *application) getCatByIDHandler(c echo.Context) error {
	cat, err := app.getCatByID(c)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, ValidationError):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, cat)
}

// Update cat godoc
//
//	@Summary		Update cat salary
//	@Description	Update cat salary by ID
//	@Tags			spycat
//	@Produce		json
//	@Param			id		path		int						true	"Cat ID"
//	@Param			payload	body		UpdateCatInfoPayload	true	"Update SpyCat payload"
//	@Success		200		{object}	store.Cat
//	@Failure		422		{object}	error
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/spycat/{id} [patch]
func (app *application) updateCatHandler(c echo.Context) error {
	var payload UpdateCatInfoPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}

	if err := Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
	}
	cat, err := app.getCatByID(c)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, ValidationError):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	updatedCat := &store.Cat{
		ID:         cat.ID,
		Name:       cat.Name,
		Experience: cat.Experience,
		Breed:      cat.Breed,
		Salary:     payload.Salary,
	}

	err = app.store.Cat.UpdateSpyCat(c.Request().Context(), updatedCat)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, updatedCat)
}

// Get cat list
//
//	@Summary		Fetches spy cat list
//	@Description	Fetches spy cat list
//	@Tags			spycat
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	[]store.Cat
//	@Failure		422		{object}	error
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/spycat [get]
func (app *application) getPaginatedCatListHandler(c echo.Context) error {
	filterDefault := store.PaginatedQuery{
		Limit:  10,
		Offset: 0,
	}
	filterQuery, err := filterDefault.Parse(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}

	if err = Validate.Struct(filterQuery); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
	}

	cats, err := app.store.Cat.GetPaginatedSpyCatList(c.Request().Context(), filterQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cats)
}

func (app *application) getCatByID(c echo.Context) (*store.Cat, error) {
	catID := c.Param("id")
	id, err := strconv.ParseInt(catID, 10, 64)
	if err != nil {
		return nil, ValidationError
	}

	cat, err := app.store.Cat.GetByID(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return nil, store.ErrNotFound
		default:
			return nil, err
		}
	}
	return cat, nil
}
