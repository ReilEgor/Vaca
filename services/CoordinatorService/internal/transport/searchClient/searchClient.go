package searchClient

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	elastic "github.com/ReilEgor/Vaca/services/SearchService/api/proto/gen/elastic"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SearchClient struct {
	client elastic.SearchServiceClient
	logger *slog.Logger
}

const (
	componentSearchClient = "searchClient"
)

func NewSearchClient(addr config.SearchClientAddr) (*SearchClient, func(), error) {
	logger := slog.With(slog.String("component", componentSearchClient))
	conn, err := grpc.NewClient(string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", domain.ErrConnectSearchService, err)
	}

	cleanup := func() {
		if err := conn.Close(); err != nil {
			logger.Error("failed to close search client connection", slog.Any("error", err))
		}
	}

	client := elastic.NewSearchServiceClient(conn)
	return &SearchClient{client: client, logger: logger}, cleanup, nil
}

func (c *SearchClient) SetVacancies(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	req := &elastic.SetVacanciesRequest{
		Result: mapModelToProto(&vacancies),
	}
	if _, err := c.client.SetVacancies(ctx, req); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrSetVacancies, err)
	}
	return nil
}

func (c *SearchClient) GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error) {
	req := &elastic.GetVacanciesRequest{
		Filter: &elastic.VacancyFilter{
			Query:        filter.Query,
			Location:     filter.Location,
			Requirements: filter.Requirements,
			Limit:        int32(filter.Limit),
			Offset:       int32(filter.Offset),
		},
	}
	resp, err := c.client.GetVacancies(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGetVacancies, err)
	}
	if resp == nil || len(resp.Vacancies) == 0 {
		return nil, nil
	}

	vacancies, err := mapProtoToModel(resp.Vacancies)
	if err != nil {
		return nil, fmt.Errorf("get vacancies: map response: %w", err)
	}
	return vacancies, nil
}

func mapModelToProto(vacancies *outPkg.ScrapeResult) *elastic.ScrapeResult {
	resultVac := make([]*elastic.Vacancy, len(vacancies.Vacancies))
	for i, v := range vacancies.Vacancies {
		resultVac[i] = &elastic.Vacancy{
			Id:           v.ID.String(),
			Title:        v.Title,
			Description:  v.Description,
			Link:         v.Link,
			Company:      v.Company,
			Salary:       v.Salary,
			Location:     v.Location,
			Requirements: v.Requirements,
			About:        v.About,
		}
	}
	if len(resultVac) > 0 {
		return &elastic.ScrapeResult{
			TaskId:    resultVac[0].Id,
			Vacancies: resultVac,
		}
	}

	return nil
}

func mapProtoToModel(vacancies []*elastic.Vacancy) ([]*outPkg.Vacancy, error) {
	result := make([]*outPkg.Vacancy, len(vacancies))
	for i, v := range vacancies {
		id, err := uuid.Parse(v.Id)
		if err != nil {
			return nil, fmt.Errorf("%w %q: %v", domain.ErrParseVacancyID, v.Id, err)
		}

		result[i] = &outPkg.Vacancy{
			ID:           id,
			Title:        v.Title,
			Description:  v.Description,
			Link:         v.Link,
			Company:      v.Company,
			Salary:       v.Salary,
			Location:     v.Location,
			Requirements: v.Requirements,
			About:        v.About,
		}
	}
	return result, nil
}
