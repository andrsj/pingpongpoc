package validator

import (
	"log/slog"

	"pingpongpoc/internal/errors"
)

type Validator struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *Validator {
	return &Validator{
		log: logger,
	}
}

func (v *Validator) ProcessInput(isValid bool) (string, error) {
	if !isValid {
		return "", errors.ErrNotValidInput
	}

	return "pong", nil
}
