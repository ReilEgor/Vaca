package domain

import "errors"

var (
	FailedToDeclareQueue     = errors.New("failed to declare queue")
	FailedToConsumeFromQueue = errors.New("failed to consume from queue")
	FailedToAckMessages      = errors.New("failed to ack messages")
	FailedToNeckMessages     = errors.New("failed to neck messages")
)
