package pkg

import "context"

type TaskPublisher interface {
	Publish(ctx context.Context, taskID string) error
}
