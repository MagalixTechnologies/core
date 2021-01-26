package utils

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// StructCtx validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified
// and also allows passing of context.Context for contextual validation information.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func StructCtx(ctx context.Context, s interface{}) interface{} {
	err := Validator.StructCtx(ctx, s)
	if err == nil {
		return nil
	}

	validationErrs := err.(validator.ValidationErrors)
	return validationErrs.Translate(Translator)
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func Struct(s interface{}) interface{} {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}

	validationErrs := err.(validator.ValidationErrors)
	return validationErrs.Translate(Translator)
}
