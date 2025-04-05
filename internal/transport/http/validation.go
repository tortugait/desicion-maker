package http

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	validationError struct {
		Field string `json:"field"`
		Tag   string `json:"tag"`
		Value any    `json:"value"`
	}

	customValidator struct {
		validator *validator.Validate
	}
)

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var ErrNotValidationErrors = errors.New("not validation errors")

func StructValidationErrors(err error, obj any) ([]validationError, error) {
	var errs []validationError
	var valErrors validator.ValidationErrors
	if errors.As(err, &valErrors) {
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		for _, vErr := range valErrors {
			field, ok := t.FieldByName(vErr.StructField())
			jsonTag := field.Tag.Get("json")

			fieldName := jsonTag
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldName = jsonTag[:idx]
			}
			if jsonTag == "" || !ok {
				fieldName = vErr.Field()
			}

			errs = append(errs, validationError{
				Field: fieldName,
				Tag:   vErr.Tag(),
				Value: vErr.Value(),
			})
		}
		return errs, nil
	}
	return nil, ErrNotValidationErrors
}
