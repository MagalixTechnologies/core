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

// StructExcept validates all fields except the ones passed in.
// Fields may be provided in a namespaced fashion relative to the  struct provided
// i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func StructExcept(s interface{}, fields ...string) interface{} {
	err := Validator.StructExcept(s, fields...)
	if err == nil {
		return nil
	}
	validationErrs := err.(validator.ValidationErrors)
	return validationErrs.Translate(Translator)
}

// StructExceptCtx validates all fields except the ones passed in and allows passing of contextual
// validation validation information via context.Context
// Fields may be provided in a namespaced fashion relative to the  struct provided
// i.e. NestedStruct.Field or NestedArrayField[0].Struct.Name
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func StructExceptCtx(ctx context.Context, s interface{}, fields ...string) interface{} {
	err := Validator.StructExceptCtx(ctx, s, fields...)
	if err == nil {
		return nil
	}
	validationErrs := err.(validator.ValidationErrors)
	return validationErrs.Translate(Translator)
}
