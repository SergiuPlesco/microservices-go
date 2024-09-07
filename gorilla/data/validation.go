package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

// ValidationError wraps the validator.FieldError to provide custom error messages
type ValidationError struct {
	Err validator.FieldError
}

func (v ValidationError) Namespace() string {
	return v.Err.Namespace()
}

func (v ValidationError) Field() string {
	return v.Err.Field()
}

func (v ValidationError) Tag() string {
	return v.Err.Tag()
}

// Error method implements the error interface
func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs {

		ve := ValidationError{err}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}

// validateSKU
func validateSKU(fl validator.FieldLevel) bool {
	// SKU must be in the format abc-abc-abc
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	sku := re.FindAllString(fl.Field().String(), -1)

	return len(sku) == 1

}
