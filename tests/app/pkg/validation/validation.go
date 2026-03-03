// Package validation provides helpers and options to configure struct validation across the app.
package validation

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10/non-standard/validators"

	uuid "bitbucket.org/csgot/helis-go-uuid"
	"github.com/go-playground/validator/v10"
)

const (
	emptyUUID = "00000000-0000-0000-0000-000000000000"
)

type (
	// ValidateOption is a func which allows to configure validate instance
	ValidateOption func(validate *validator.Validate)
)

// NewStructValidator returns *validator.Validate with all custom functions registered
func NewStructValidator(opts ...ValidateOption) (*validator.Validate, error) {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(uuidValuer, uuid.UUID{})

	err := validate.RegisterValidation("not_blank", validators.NotBlank)
	if err != nil {
		return nil, fmt.Errorf("failed to register not_blank validation: %w", err)
	}

	for _, opt := range opts {
		opt(validate)
	}

	return validate, nil
}

func uuidValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(uuid.UUID); ok {
		strVal := valuer.String()
		if strVal == emptyUUID {
			return nil
		}

		return strVal
	}

	return nil
}
