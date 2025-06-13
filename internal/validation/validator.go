package validation

import (
	"github.com/go-playground/validator/v10"
)

func RegisterCatValidator(v *validator.Validate, client *Client) {
	v.RegisterValidation("breed-exits", func(fl validator.FieldLevel) bool {
		breed := fl.Field().String()
		ok, err := client.CatBreedExists(breed)
		return err == nil && ok
	})
}
