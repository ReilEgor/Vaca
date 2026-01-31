package elasticsearch

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ElasticRepository struct {
	client *elastic.TypedClient
	logger *slog.Logger
}

func NewElasticRepository(client *elastic.TypedClient) domain.VacancySearchRepository {
	return &ElasticRepository{client: client, logger: slog.With(slog.String("component", "ElasticRepository"))}
}

func (e *ElasticRepository) Index(ctx context.Context, v outPkg.Vacancy, taskID string) error {
	document := struct {
		TaskID       string `json:"task_id"`
		Title        string `json:"title"`
		Company      string `json:"company"`
		Location     string `json:"location"`
		Salary       string `json:"salary"`
		Description  string `json:"description"`
		About        string `json:"about"`
		Requirements string `json:"requirements"`
		Link         string `json:"link"`
	}{
		TaskID:       taskID,
		Title:        v.Title,
		Company:      v.Company,
		Location:     v.Location,
		Salary:       v.Salary,
		Description:  v.Description,
		About:        v.About,
		Requirements: v.Requirements,
		Link:         v.Link,
	}
	_, err := e.client.Index("vacancies").
		Id(taskID).
		Request(document).
		Do(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticRepository) IndexBatch(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	if len(vacancies.Vacancies) == 0 {
		e.logger.Info("no vacancies to index, skipping bulk request", slog.String("task_id", vacancies.TaskID.String()))
		return nil
	}
	bulk := e.client.Bulk().Index("vacancies")

	for _, v := range vacancies.Vacancies {
		doc := struct {
			TaskID       string `json:"task_id"`
			Title        string `json:"title"`
			Company      string `json:"company"`
			Location     string `json:"location"`
			Salary       string `json:"salary"`
			Description  string `json:"description"`
			About        string `json:"about"`
			Requirements string `json:"requirements"`
			Link         string `json:"link"`
		}{
			TaskID:       vacancies.TaskID.String(),
			Title:        v.Title,
			Company:      v.Company,
			Location:     v.Location,
			Salary:       v.Salary,
			Description:  v.Description,
			About:        v.About,
			Requirements: v.Requirements,
			Link:         v.Link,
		}

		err := bulk.IndexOp(types.IndexOperation{Id_: &v.Link}, doc)
		if err != nil {
			e.logger.Error("failed to add vacancy to bulk operation",
				slog.String("url", v.Link),
				slog.Any("err", err))
			continue
		}
	}

	res, err := bulk.Do(ctx)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}

	if res.Errors {
		e.logger.Warn("bulk operation finished with some errors")
	}

	e.logger.Info("successfully indexed batch of vacancies",
		slog.Int("count", len(vacancies.Vacancies)),
		slog.String("task_id", vacancies.TaskID.String()))

	return nil
}

func (e ElasticRepository) Search(ctx context.Context, query string) ([]outPkg.Vacancy, error) {
	return nil, nil
}
