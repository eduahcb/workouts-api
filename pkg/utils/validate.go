package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidateMessageError []map[string]string

func transforToTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// skip if tag key says it should be ignored
	if name == "-" {
		return ""
	}
	return name
}

func ValidateStruct(s any) (ValidateMessageError, error) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(transforToTagName)

	err := validate.Struct(s)

	if err != nil {
		var messages []map[string]string

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			element := map[string]string{
				"field":   e.Field(),
				"message": e.Tag(),
			}

			messages = append(messages, element)
		}

		return messages, errors.New("validation failed")
	}

	return nil, nil
}
