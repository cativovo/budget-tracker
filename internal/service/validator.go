package service

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New(validator.WithRequiredStructEnabled())

	// https://github.com/go-playground/validator/issues/861
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}

func validateStruct(s any) error {
	if err := v.Struct(s); err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return errors.New("invalid ValidationErrors")
		}

		m := make([]string, 0, len(vErrs))

		for _, err := range vErrs {
			m = append(m, getMessage(s, err))
		}

		return errors.New(strings.Join(m, ", "))
	}
	return nil
}

func getMessage(v any, err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("invalid %s", err.Field())
	case "required_with":
		if field, ok := reflect.TypeOf(v).Elem().FieldByName(err.Param()); ok {
			if jsonTag, ok := field.Tag.Lookup("json"); ok {
				return fmt.Sprintf("%s is required with %s", err.Field(), jsonTag)
			}
		}
		return fmt.Sprintf("%s is required", err.Field())
	case "number":
		return fmt.Sprintf("%s must have a valid numeric value", err.Field())
	case "hexcolor":
		return fmt.Sprintf("%s must have a valid hex color value", err.Field())
	case "datetime":
		return fmt.Sprintf("invalid %s", err.Field())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
	case "min":
		p := err.Param()
		// https://github.com/go-playground/validator/issues/1308
		if p == "1" {
			return fmt.Sprintf("%s is required", err.Field())
		}

		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s: '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
	}
}
