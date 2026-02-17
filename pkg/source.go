package pkg

import (
	"context"

	taskstate "github.com/ReilEgor/Vaca/services/StateService/api/proto/gen/task_state"
	"github.com/google/uuid"
)

type Source struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	IsActive bool      `json:"is_active"`
}

func MapDomainSourceToProto(s Source) *taskstate.Source {
	return &taskstate.Source{
		Id:       s.ID.String(),
		Name:     s.Name,
		Url:      s.URL,
		IsActive: s.IsActive,
	}
}

func MapProtoToDomainSource(s *taskstate.Source) Source {
	id, err := uuid.Parse(s.Id)
	if err != nil {
		return Source{}
	}
	return Source{
		ID:       id,
		Name:     s.Name,
		URL:      s.Url,
		IsActive: s.IsActive,
	}
}

type SourceRepository interface {
	GetAllActiveSources(ctx context.Context) ([]*Source, error)
}
