package pingpong

import (
	"context"
	"fmt"
	"log/slog"
)

type (
	InputFilter interface {
		Filter(context.Context, string) (bool, error)
	}

	Validator interface {
		ProcessInput(bool) (string, error)
	}

	Service struct {
		filter    InputFilter
		validator Validator

		log *slog.Logger
	}
)

func NewService(filter InputFilter, validator Validator, logger *slog.Logger) *Service {
	return &Service{
		filter:    filter,
		log:       logger,
		validator: validator,
	}
}

func (s *Service) Serve(ctx context.Context, input string) (string, error) {
	s.log.Info("Executing service 'PingPong'", "input", input)

	isFiltered, err := s.filter.Filter(ctx, input)
	if err != nil {
		return "", fmt.Errorf("filtering error: %w", err)
	}

	result, err := s.validator.ProcessInput(isFiltered)
	if err != nil {
		return "", fmt.Errorf("error validate input: %w", err)
	}

	return result, nil
}
