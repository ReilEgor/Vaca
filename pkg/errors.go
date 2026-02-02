package pkg

import "errors"

var (
	ErrTaskIDRequired   = errors.New("task ID is required")
	ErrSourceIDRequired = errors.New("source ID is required")

	ErrInvalidRequest = errors.New("invalid request")
)
