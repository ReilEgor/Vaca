package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	taskstate "github.com/ReilEgor/Vaca/services/StateService/api/proto/gen/task_state"
	"github.com/ReilEgor/Vaca/services/StateService/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskStateHandler struct {
	taskstate.UnimplementedTaskStateServiceServer
	usecase domain.StateUsecase
	logger  *slog.Logger
}

func NewTaskStateHandler(uc domain.StateUsecase) *TaskStateHandler {
	return &TaskStateHandler{
		usecase: uc,
		logger:  slog.With(slog.String("component", "taskStateHandler")),
	}
}

func (h *TaskStateHandler) SetTask(ctx context.Context, req *taskstate.SetTaskRequest) (*taskstate.SetTaskResponse, error) {
	h.logger.Debug("SetTask called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err := h.usecase.Set(ctx, req.TaskId, req.SearchKey, int(req.TotalSources), time.Minute*2)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "request timed out")
		}
		return nil, err
	}
	resp := &taskstate.SetTaskResponse{}
	return resp, nil
}

func (h *TaskStateHandler) GetTask(ctx context.Context, req *taskstate.GetTaskRequest) (*taskstate.GetTaskResponse, error) {
	h.logger.Debug("GetTask called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	task := h.usecase.Get(ctx, req.TaskId)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	}
	resp := &taskstate.GetTaskResponse{
		Values: task,
	}
	return resp, nil
}

func (h *TaskStateHandler) GetTaskIDByHash(ctx context.Context, req *taskstate.GetTaskIDByHashRequest) (*taskstate.GetTaskIDByHashResponse, error) {
	h.logger.Debug("GetTaskIDByHash called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	id, err := h.usecase.GetIDByHash(ctx, req.SearchKey)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "request timed out")
		}
		return nil, err
	}
	resp := &taskstate.GetTaskIDByHashResponse{
		Id: id,
	}
	return resp, nil
}

func (h *TaskStateHandler) GetSources(ctx context.Context, req *taskstate.GetSourcesRequest) (*taskstate.GetSourcesResponse, error) {
	h.logger.Debug("GetSources called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	sources, err := h.usecase.GetSources(ctx)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "request timed out")
		}
		return nil, err
	}
	respSources := make([]*taskstate.Source, len(sources))
	for i, source := range sources {
		respSources[i] = outPkg.MapDomainSourceToProto(source)
	}
	resp := &taskstate.GetSourcesResponse{
		Sources: respSources,
	}
	return resp, nil
}

func (h *TaskStateHandler) SetStatus(ctx context.Context, req *taskstate.SetStatusRequest) (*taskstate.SetStatusResponse, error) {
	h.logger.Debug("SetStatus called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err := h.usecase.SetStatus(ctx, req.TaskId, req.Status)
	if err != nil {
		return nil, err
	}
	resp := &taskstate.SetStatusResponse{}
	return resp, nil
}

func (h *TaskStateHandler) GetTotal(ctx context.Context, req *taskstate.GetTotalRequest) (*taskstate.GetTotalResponse, error) {
	h.logger.Debug("GetTotal called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	total, err := h.usecase.GetTotal(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}
	resp := &taskstate.GetTotalResponse{
		Total: int32(total),
	}
	return resp, nil
}

func (h *TaskStateHandler) IncrementCompleted(ctx context.Context, req *taskstate.IncrementCompletedRequest) (*taskstate.IncrementCompletedResponse, error) {
	h.logger.Debug("GetTotal called")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	total, err := h.usecase.IncrementCompleted(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}
	resp := &taskstate.IncrementCompletedResponse{
		Total: int32(total),
	}
	return resp, nil
}

func (h *TaskStateHandler) Register(ctx context.Context, req *taskstate.RegisterRequest) (*taskstate.RegisterResponse, error) {
	id, err := uuid.Parse(req.Source.Id)
	if err != nil {
		return nil, err
	}
	source := outPkg.Source{
		ID:       id,
		Name:     req.Source.Name,
		URL:      req.Source.Url,
		IsActive: req.Source.IsActive,
	}
	err = h.usecase.Register(ctx, source)
	if err != nil {
		return nil, err
	}
	resp := &taskstate.RegisterResponse{}
	return resp, nil
}
