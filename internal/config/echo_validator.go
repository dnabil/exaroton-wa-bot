package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IValidator interface {
	Validate(any) error
}

type Validator struct{}

func (*Validator) Validate(obj any) error {
	validatable, ok := obj.(validation.Validatable)

	if ok {
		return validatable.Validate()
	}

	return nil
}

func NewValidator() IValidator {
	return &Validator{}
}
