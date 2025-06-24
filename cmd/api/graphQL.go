package main

// Get one Cat
//
//	@Summary		List of cats
//	@Description	List of cats
//	@Tags			spycat
//	@Param			id	path		int	true	"Cat ID"
//	@Success		200	{object}	string
//	@Failure		500	{object}	error
//	@Router			/ql/{id} [get]
/*func (app *application) getListOfCatsQL(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError)
	}

	data := app.graphqlStorage.Cat.GetOneCat(c.Request().Context(), idInt)
	if len(data.Errors) > 0 {
		return c.JSON(http.StatusInternalServerError, data.Errors)
	}

	return c.JSON(http.StatusOK, data)
}*/
