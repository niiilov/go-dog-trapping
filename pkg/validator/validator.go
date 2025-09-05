package validate

import (
	"github.com/go-playground/validator/v10"
)

func Validate(v any) error {
	valid := validator.New()

	return valid.Struct(v)
}
