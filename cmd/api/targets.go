package main

import (
	"FIDOtestBackendApp/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Target struct {
	Name     string `json:"name" validate:"required,max=200,min=1"`
	Country  string `json:"country" validate:"required,max=200,min=1"`
	Notes    string `json:"notes" validate:"required,max=255,min=1"`
	Complete *bool  `json:"complete" validate:"required"`
}

type UpdateNotesPayload struct {
	Notes string `json:"notes" validate:"required,max=255,min=1"`
}

// Update target's note
//
//	@Summary		Update target's note
//	@Description	Update target's note  by ID
//	@Tags			target
//	@Produce		json
//	@Param			mission_id	path		int					true	"mission_id's ID"
//	@Param			target_id	path		int					true	"target_id's ID"
//	@Param			payload		body		UpdateNotesPayload	true	"Update Target note"
//	@Success		200			{object}	store.UpdateTargetNote
//	@Failure		422			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Router			/mission/{mission_id}/target/{target_id} [patch]
func (app *application) updateTargetNote(c echo.Context) error {
	var payload UpdateNotesPayload
	parsedNoteId, parsedMissionId, err := parseParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}
	if err = Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
	}
	updateNote := &store.UpdateTargetNote{
		ID:        parsedNoteId,
		MissionID: parsedMissionId,
		Note:      payload.Notes,
	}
	err = app.store.Target.UpdateTargetNote(c.Request().Context(), updateNote)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, updateNote)
}

// Update target's status
//
//	@Summary		Update target's status
//	@Description	Update target's status  by ID
//	@Tags			target
//	@Produce		json
//	@Param			mission_id	path		int	true	"mission_id's ID"
//	@Param			target_id	path		int	true	"target_id's ID"
//	@Success		200			{object}	store.UpdateTargetStatus
//	@Failure		422			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Router			/mission/{mission_id}/target_status/{target_id} [patch]
func (app *application) updateTargetStatus(c echo.Context) error {
	parsedNoteId, parsedMissionId, err := parseParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	updateNote := &store.UpdateTargetStatus{
		ID:        parsedNoteId,
		MissionID: parsedMissionId,
		Status:    true,
	}
	err = app.store.Target.UpdateTargetStatus(c.Request().Context(), updateNote)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, updateNote)
}

// Delete target
//
//	@Summary		Delete target
//	@Description	Delete target by mission id and target id
//	@Tags			target
//	@Produce		json
//	@Param			mission_id	path		int	true	"mission_id's ID"
//	@Param			target_id	path		int	true	"target_id's ID"
//	@Success		204			{object}	nil
//	@Failure		422			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Router			/mission/{mission_id}/target/{target_id} [delete]
func (app *application) deleteTarget(c echo.Context) error {
	parsedNoteId, parsedMissionId, err := parseParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = app.store.Target.DeleteTarget(c.Request().Context(), parsedMissionId, parsedNoteId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// Add target
//
//	@Summary		Add target to mission
//	@Description	Add target to mission by mission_id and target_id
//	@Tags			target
//	@Produce		json
//	@Param			mission_id	path		int		true	"mission_id's ID"
//	@Param			payload		body		Target	true	"Target payload"
//	@Success		204			{object}	nil
//	@Failure		422			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Router			/mission/{mission_id}/target [post]
func (app *application) addTarget(c echo.Context) error {
	var payload Target
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}
	missionId := c.Param("mission_id")
	parsedMissionId, err := strconv.ParseInt(missionId, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError.Error())
	}
	if err = Validate.Struct(payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ValidationError.Error())
	}
	target := &store.Target{
		MissionID: parsedMissionId,
		Name:      payload.Name,
		Country:   payload.Country,
		Notes:     payload.Notes,
	}
	err = app.store.Target.AddTarget(c.Request().Context(), target)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return c.JSON(http.StatusUnprocessableEntity, err.Error())
		case store.TargetAmountError:
			return c.JSON(http.StatusConflict, err.Error())
		case store.ViolatePK:
			return c.JSON(http.StatusConflict, err.Error())
		default:
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusCreated)
}

func parseParams(c echo.Context) (targetId, missionId int64, err error) {
	parsedMissionId, err := strconv.ParseInt(c.Param("mission_id"), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	parsedTargetId, err := strconv.ParseInt(c.Param("target_id"), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	return parsedTargetId, parsedMissionId, nil
}
