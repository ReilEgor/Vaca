package stateClient

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	taskstate "github.com/ReilEgor/Vaca/services/StateService/api/proto/gen/task_state"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StateClient struct {
	client taskstate.TaskStateServiceClient
	logger *slog.Logger
}

const (
	componentStateClient = "stateClient"
)

func NewStateClient(addr config.StateClientAddr) (*StateClient, func(), error) {
	logger := slog.With(slog.String("component", componentStateClient))
	conn, err := grpc.NewClient(string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", domain.ErrConnectStateService, err)
	}
	cleanup := func() {
		if err := conn.Close(); err != nil {
			logger.Error("failed to close state client connection", slog.Any("error", err))
		}
	}

	client := taskstate.NewTaskStateServiceClient(conn)
	return &StateClient{client: client, logger: logger}, cleanup, nil
}

func (c *StateClient) Set(ctx context.Context, taskID string, searchKey string, totalSources int) error {
	req := &taskstate.SetTaskRequest{
		TaskId:       taskID,
		SearchKey:    searchKey,
		TotalSources: int32(totalSources),
	}

	if _, err := c.client.SetTask(ctx, req); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrSetTask, err)
	}

	c.logger.Info("set task completed", slog.String("task_id", taskID))

	return nil
}

func (c *StateClient) Get(ctx context.Context, taskID string) (map[string]string, error) {
	req := &taskstate.GetTaskRequest{
		TaskId: taskID,
	}
	resp, err := c.client.GetTask(ctx, req) // TODO: validate
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGetTask, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("%w: get task %q", domain.ErrNilResponse, taskID)
	}

	c.logger.Info("get task completed", slog.String("task_id", taskID))

	return resp.Values, nil
}

func (c *StateClient) GetIDByHash(ctx context.Context, searchKey string) (string, error) {
	req := &taskstate.GetTaskIDByHashRequest{
		SearchKey: searchKey,
	}
	resp, err := c.client.GetTaskIDByHash(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrGetTaskIDByHash, err)
	}
	if resp == nil {
		return "", fmt.Errorf("%w: get task id by hash %q", domain.ErrNilResponse, searchKey)
	}

	c.logger.Info("get id by hash completed", slog.String("search_key", searchKey))

	return resp.Id, nil
}

func (c *StateClient) GetSources(ctx context.Context) ([]outPkg.Source, error) {
	resp, err := c.client.GetSources(ctx, &taskstate.GetSourcesRequest{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGetSources, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("%w: get sources", domain.ErrNilResponse)
	}
	respSources := make([]outPkg.Source, len(resp.Sources))

	for i, source := range resp.Sources {
		respSources[i] = outPkg.MapProtoToDomainSource(source)
	}
	c.logger.Info("get sources completed", slog.Int("count", len(respSources)))

	return respSources, nil
}

func (c *StateClient) IncrementCompleted(ctx context.Context, taskID string) (int64, error) {
	req := &taskstate.IncrementCompletedRequest{
		TaskId: taskID,
	}
	resp, err := c.client.IncrementCompleted(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", domain.ErrIncrementCompleted, err)
	}
	if resp == nil {
		return 0, fmt.Errorf("%w: increment completed %q", domain.ErrNilResponse, taskID)
	}
	return int64(resp.Total), nil
}
