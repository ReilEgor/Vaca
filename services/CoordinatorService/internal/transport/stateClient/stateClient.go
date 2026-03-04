package stateClient

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	taskstate "github.com/ReilEgor/Vaca/services/StateService/api/proto/gen/task_state"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StateClient struct {
	client taskstate.TaskStateServiceClient
	logger *slog.Logger
}

func NewStateClient(addr config.StateClientAddr) *StateClient {
	conn, err := grpc.NewClient(string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	logger := slog.With(slog.String("component", "redisStatusRepository"))
	client := taskstate.NewTaskStateServiceClient(conn)
	return &StateClient{client: client, logger: logger}
}

func (c *StateClient) Set(ctx context.Context, taskID string, searchKey string, totalSources int) error {
	req := &taskstate.SetTaskRequest{
		TaskId:       taskID,
		SearchKey:    searchKey,
		TotalSources: int32(totalSources),
	}

	_, err := c.client.SetTask(ctx, req)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	c.logger.Info("Set task completed", taskID)
	return nil
}

func (c *StateClient) Get(ctx context.Context, taskID string) map[string]string {
	req := &taskstate.GetTaskRequest{
		TaskId: taskID,
	}
	resp, err := c.client.GetTask(ctx, req) // TODO: validate
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}
	c.logger.Info("Get task completed")
	return resp.Values
}

func (c *StateClient) GetIDByHash(ctx context.Context, searchKey string) (string, error) {
	req := &taskstate.GetTaskIDByHashRequest{
		SearchKey: searchKey,
	}
	resp, err := c.client.GetTaskIDByHash(ctx, req)
	c.logger.Info("GetIDByHash completed")
	return resp.Id, err
}

func (c *StateClient) GetSources(ctx context.Context) ([]outPkg.Source, error) {
	req := &taskstate.GetSourcesRequest{}
	resp, err := c.client.GetSources(ctx, req)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}
	respSources := make([]outPkg.Source, len(resp.Sources))

	for i, source := range resp.Sources {
		respSources[i] = outPkg.MapProtoToDomainSource(source)
	}
	c.logger.Info("GetSources completed")
	return respSources, nil
}

func (c *StateClient) IncrementCompleted(ctx context.Context, taskID string) (int64, error) {
	req := &taskstate.IncrementCompletedRequest{
		TaskId: taskID,
	}
	resp, err := c.client.IncrementCompleted(ctx, req)
	if err != nil {
		return 0, err
	}
	return int64(resp.Total), nil
}
