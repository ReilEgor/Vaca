package domain

import (
	"errors"
)

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrFailedToCreateTask = errors.New("failed to create task")
	ErrTaskAlreadyExists  = errors.New("task already exists")

	ErrFailedToGetSources = errors.New("failed to get sources")

	ErrFailedToGetVacancies = errors.New("failed to get vacancies")

	ErrInvalidSource     = errors.New("invalid or unsupported source")
	ErrSourceUnavailable = errors.New("source is temporarily unavailable")

	ErrInvalidKeywords = errors.New("keywords cannot be empty")

	ErrInvalidRequestBody = errors.New("invalid request body")
)
