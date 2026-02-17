package domain

import "errors"

var (
	// Redis errors
	FailedToIncrementCompleted = errors.New("failed to increment completed")
	FailedToUpdateStatus       = errors.New("failed to update status")
	FailedToGetTotal           = errors.New("failed to get total")
)
