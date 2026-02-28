package domain

import (
	"errors"
)

var (
	// --- Task Management Errors ---
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskAlreadyExists  = errors.New("task already exists")
	ErrFailedToCreateTask = errors.New("failed to create task")
	ErrSearchFailed       = errors.New("search failed")
	ErrFailedToGetTask    = errors.New("failed to get task status")
	ErrInvalidTaskID      = errors.New("invalid task ID")
	
	// --- Source & Content Errors ---
	ErrFailedToGetSources   = errors.New("failed to retrieve sources")
	ErrFailedToGetVacancies = errors.New("failed to retrieve vacancies")
	ErrInvalidSource        = errors.New("invalid or unsupported source")
	ErrSourceUnavailable    = errors.New("source is temporarily unavailable")

	// --- Validation & Data Errors ---
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidKeywords    = errors.New("keywords cannot be empty")
	ErrMarshalTask        = errors.New("failed to marshal task message")

	// --- Infrastructure: RabbitMQ Lifecycle ---
	ErrConnect      = errors.New("failed to connect to rabbitmq")
	ErrOpenChannel  = errors.New("failed to open rabbitmq channel")
	ErrCloseConn    = errors.New("failed to close rabbitmq connection")
	ErrCloseChannel = errors.New("failed to close rabbitmq channel")

	// --- Infrastructure: RabbitMQ Operations ---
	ErrDeclareExchange = errors.New("failed to declare exchange")
	ErrPublishMessage  = errors.New("failed to publish message")

	// --- Transport: Gin Errors ---
	ErrServerNotInitialized = errors.New("server is not initialized")

	// --- Transport: Search Service Client Errors ---
	ErrConnectSearchService = errors.New("connect to search service")
	ErrSetVacancies         = errors.New("set vacancies")
	ErrGetVacancies         = errors.New("get vacancies")
	ErrParseVacancyID       = errors.New("parse vacancy id")

	// --- Transport: State Service Client Errors ---
	ErrConnectStateService = errors.New("connect to state service")
	ErrSetTask             = errors.New("set task")
	ErrGetTask             = errors.New("get task")
	ErrGetTaskIDByHash     = errors.New("get task id by hash")
	ErrGetSources          = errors.New("get sources")
	ErrIncrementCompleted  = errors.New("increment completed")
	ErrNilResponse         = errors.New("nil response from state service")
)
