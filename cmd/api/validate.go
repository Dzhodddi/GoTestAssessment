package main

import (
	"FIDOtestBackendApp/internal/validation"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	catValidator := validation.NewClient("https://api.thecatapi.com/v1/breeds")
	validation.RegisterCatValidator(Validate, catValidator)
}
