package models

import "github.com/go-playground/validator/v10"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationError(err validator.FieldError) ValidationError {
	return ValidationError{
		Field:   err.Field(),
		Message: getValidationErrorMessage(err),
	}
}

func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		if err.Type().String() == "[]interface {}" {
			return "At least one item is required"
		}
		return "This field must be at least " + err.Param() + " characters long"
	default:
		return "Invalid value"
	}
}
