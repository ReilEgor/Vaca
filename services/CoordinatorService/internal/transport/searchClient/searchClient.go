package searchClient

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/config"
	elastic "github.com/ReilEgor/Vaca/services/SearchService/api/proto/gen/elastic"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SearchClient struct {
	client elastic.SearchServiceClient
	logger *slog.Logger
}

func NewSearchClient(addr config.SearchClientAddr) *SearchClient {
	conn, err := grpc.NewClient(string(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	logger := slog.With(slog.String("component", "redisStatusRepository"))
	client := elastic.NewSearchServiceClient(conn)
	return &SearchClient{client: client, logger: logger}
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

func mapProtoToModel(vacancies []*elastic.Vacancy) []*outPkg.Vacancy {
	resultVac := make([]*outPkg.Vacancy, len(vacancies))
	for i, v := range vacancies {
		id, err := uuid.Parse(v.Id)
		if err != nil {
			return nil
		}
		resultVac[i] = &outPkg.Vacancy{
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
	if len(resultVac) > 0 {
		return resultVac
	}

	return nil
}

func (c *SearchClient) SetVacancies(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	req := elastic.SetVacanciesRequest{
		Result: mapModelToProto(&vacancies),
	}
	_, err := c.client.SetVacancies(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}

func (c *SearchClient) GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error) {
	req := elastic.GetVacanciesRequest{
		Filter: &elastic.VacancyFilter{
			Query:        filter.Query,
			Location:     filter.Location,
			Requirements: filter.Requirements,
			Limit:        int32(filter.Limit),
			Offset:       int32(filter.Offset),
		},
	}
	resp, err := c.client.GetVacancies(ctx, &req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	vacancies := mapProtoToModel(resp.Vacancies)
	return vacancies, nil
}
