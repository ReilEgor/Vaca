package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/operator"
)

type ElasticRepository struct {
	client *elastic.TypedClient
	logger *slog.Logger
}

func NewElasticRepository(client *elastic.TypedClient) *ElasticRepository {
	return &ElasticRepository{client: client, logger: slog.With(slog.String("component", "elasticRepository"))}
}

func (r *ElasticRepository) Search(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error) {
	var mustConditions []types.Query

	if filter.Query != "" {
		mustConditions = append(mustConditions, types.Query{
			Match: map[string]types.MatchQuery{
				"description": {
					Query:     filter.Query,
					Operator:  &operator.And,
					Fuzziness: "AUTO",
				},
			},
		})
	}

	if filter.Location != "" {
		mustConditions = append(mustConditions, types.Query{
			Match: map[string]types.MatchQuery{
				"location": {
					Query:     filter.Location,
					Fuzziness: "AUTO",
				},
			},
		})
	}

	searchQuery := &types.Query{
		Bool: &types.BoolQuery{
			Must: mustConditions,
		},
	}

	res, err := r.client.Search().
		Index("vacancies").
		Query(searchQuery).
		From(filter.Offset).
		Size(filter.Limit).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	return r.mapHitsToVacancies(res.Hits.Hits), nil
}

func (r *ElasticRepository) mapHitsToVacancies(hits []types.Hit) []*outPkg.Vacancy {
	vacancies := make([]*outPkg.Vacancy, 0, len(hits))

	for _, hit := range hits {
		v := new(outPkg.Vacancy)
		if err := json.Unmarshal(hit.Source_, v); err != nil {
			r.logger.Error("unmarshal error", slog.Any("error", err))
			continue
		}

		vacancies = append(vacancies, v)
	}
	return vacancies
}
