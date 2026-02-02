package domain

import "errors"

var (
	FailedToDeclareQueue     = errors.New("failed to declare queue")
	FailedToConsumeFromQueue = errors.New("failed to consume from queue")
	FailedToNeckMessages     = errors.New("failed to neck messages")
	FailedToAckMessages      = errors.New("failed to ack messages")

	FailedBulkRequest = errors.New("failed bulk request")

	//Redis errors
	FailedToIncrementCompleted = errors.New("failed to increment completed")
	FailedToUpdateStatus       = errors.New("failed to update status")
	FailedToGetTotal           = errors.New("failed to get total")
)
