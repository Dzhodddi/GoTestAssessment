package main

import (
	"FIDOtestBackendApp/internal/store"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type MissionPayload struct {
	Complete bool     `json:"complete" validate:"required"`
	Targets  []Target `json:"targets" validate:"required,dive"`
}

// Create Mission
//
//	@Summary		Create Mission
//	@Description	Create new Mission
//	@Tags			mission
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		MissionPayload	true	"Mission payload"
//	@Success		201		{object}	store.MissionWithTargets
//	@Failure		400		{object}	error
//	@Failure		422		{object}	error
//	@Failure		500		{object}	error
//	@Router			/mission [post]
func (app *application) createMissionHandler(c echo.Context) error {
	var payload MissionPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}

	if err := Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
	}

	targets := make([]store.Target, 0, len(payload.Targets))
	nameSet := make(map[string]struct{})

	for _, target := range payload.Targets {
		if _, exists := nameSet[target.Name]; exists {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("duplicate target name: %s", target.Name))
		}
		nameSet[target.Name] = struct{}{}

		targets = append(targets, store.Target{
			Name:      target.Name,
			Country:   target.Country,
			Notes:     target.Notes,
			Completed: target.Complete,
		})
	}

	mission := &store.MissionWithTargets{
		Targets: targets,
		Mission: store.Mission{
			CatID:     nil,
			Completed: payload.Complete,
		},
	}

	err := app.store.Mission.CreateMission(c.Request().Context(), mission)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return nil
}

// Delete mission
//
//	@Summary		Delete mission
//	@Description	Delete mission by ID
//	@Tags			mission
//	@Produce		json
//	@Param			id	path		int	true	"mission ID"
//	@Success		204	{object}	nil
//	@Failure		422	{object}	error
//	@Failure		400	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/mission/{id} [delete]
func (app *application) deleteMissionHandler(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}

	err = app.store.Mission.DeleteMission(c.Request().Context(), parsedID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, store.MissionedAssigned):
			return c.JSON(http.StatusConflict, ConflictError.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func (app *application) getMissionById(c echo.Context) error {
	return c.JSON(http.StatusOK, app.store.Mission)
}

// Update mission
//
//	@Summary		Update mission
//	@Description	Update mission by ID
//	@Tags			mission
//	@Produce		json
//	@Param			id	path		int	true	"Cat ID"
//	@Success		200	{object}	store.Cat
//	@Failure		422	{object}	error
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/mission/{id} [patch]
func (app *application) updateMissionStatus(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
	}
	payload := &store.UpdatedMission{
		ID:     parsedID,
		Status: true,
	}
	err = app.store.Mission.UpdateMissionStatus(c.Request().Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, payload)
}
