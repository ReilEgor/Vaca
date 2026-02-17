package domain

import (
	"context"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	taskstate "github.com/ReilEgor/Vaca/services/StateService/api/proto/gen/task_state"
)

//go:generate mockery --name StateRepository --output ../mocks/domain --outpkg domain --case=underscore
type StateRepository interface {
	Set(ctx context.Context, taskID string, searchKey string, totalSources int, ttl time.Duration) error
	Get(ctx context.Context, taskID string) map[string]string
	GetIDByHash(ctx context.Context, searchKey string) (string, error)
	GetSources(ctx context.Context) ([]outPkg.Source, error)
	SetStatus(ctx context.Context, taskID string, status string) error
	GetTotal(ctx context.Context, taskID string) (int64, error)
	IncrementCompleted(ctx context.Context, taskID string) (int64, error)
	Register(ctx context.Context, source outPkg.Source) error
}

type StateUsecase interface {
	Set(ctx context.Context, taskID string, searchKey string, totalSources int, ttl time.Duration) error
	Get(ctx context.Context, taskID string) map[string]string
	GetIDByHash(ctx context.Context, searchKey string) (string, error)
	GetSources(ctx context.Context) ([]outPkg.Source, error)
	SetStatus(ctx context.Context, taskID string, status string) error
	GetTotal(ctx context.Context, taskID string) (int64, error)
	IncrementCompleted(ctx context.Context, taskID string) (int64, error)
	Register(ctx context.Context, source outPkg.Source) error
}

type StateHandler interface {
	SetTask(ctx context.Context, req *taskstate.SetTaskRequest) (*taskstate.SetTaskResponse, error)
	GetTask(ctx context.Context, req *taskstate.GetTaskRequest) (*taskstate.GetTaskResponse, error)
	GetTaskIDByHash(ctx context.Context, req *taskstate.GetTaskIDByHashRequest) (*taskstate.GetTaskIDByHashResponse, error)
	GetSources(ctx context.Context, req *taskstate.GetSourcesRequest) (*taskstate.GetSourcesResponse, error)
	SetStatus(ctx context.Context, req *taskstate.SetStatusRequest) (*taskstate.SetStatusResponse, error)
	GetTotal(ctx context.Context, req *taskstate.GetTotalRequest) (*taskstate.GetTotalResponse, error)
	IncrementCompleted(ctx context.Context, req *taskstate.IncrementCompletedRequest) (*taskstate.IncrementCompletedResponse, error)
	Register(ctx context.Context, req *taskstate.RegisterRequest) (*taskstate.RegisterResponse, error)
}
