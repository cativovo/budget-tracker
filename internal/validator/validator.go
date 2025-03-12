package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/go-playground/validator/v10"
)

// https://github.com/go-playground/validator/issues/559#issuecomment-1871786235
type validationErrors []error

func (v validationErrors) Error() string {
	var message string

	for i, err := range v {
		if i > 0 {
			message += ","
		}
		message += err.Error()
	}

	return message
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Struct(s any) error {
	if err := v.validator.Struct(s); err != nil {
		var vErrs validationErrors

		for _, err := range err.(validator.ValidationErrors) {
			var e error
			switch err.Tag() {
			case "required":
				e = fmt.Errorf("'%s' is required", err.Field())
			case "required_with":
				if field, ok := reflect.TypeOf(s).Elem().FieldByName(err.Param()); ok {
					if jsonTag, ok := field.Tag.Lookup("json"); ok {
						e = fmt.Errorf("'%s' is required with '%s'", err.Field(), jsonTag)
					}
				}
			case "number":
				e = fmt.Errorf("'%s' must have a valid numeric value", err.Field())
			case "hexcolor":
				e = fmt.Errorf("'%s' must have a valid hex color value", err.Field())
			case "datetime":
				e = fmt.Errorf("'%s' must have a valid date value", err.Field())
			case "gte":
				e = fmt.Errorf("'%s' must be greater than or equal to %s", err.Field(), err.Param())
			case "gt":
				e = fmt.Errorf("'%s' must be greater than %s", err.Field(), err.Param())
			default:
				e = fmt.Errorf("'%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
			}
			vErrs = append(vErrs, e)
		}

		return internal.NewError(internal.ErrorCodeInvalid, vErrs.Error())
	}

	return nil
}

func (v *Validator) Var(f any, tag string) error {
	return v.validator.Var(f, tag)
}

func NewValidator() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	// https://github.com/go-playground/validator/issues/861
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: v,
	}
}
