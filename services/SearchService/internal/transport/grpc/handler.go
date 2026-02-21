package grpc

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	search "github.com/ReilEgor/Vaca/services/SearchService/api/proto/gen/elastic"
	"github.com/ReilEgor/Vaca/services/SearchService/internal/domain"
	"github.com/google/uuid"
)

type SearchHandler struct {
	search.UnimplementedSearchServiceServer
	usecase domain.SearchUsecase
	logger  *slog.Logger
}

func NewSearchHandler(uc domain.SearchUsecase) *SearchHandler {
	return &SearchHandler{
		usecase: uc,
		logger:  slog.With(slog.String("component", "searchHandler")),
	}
}

func mapProtoToModel(p *search.Vacancy) *outPkg.Vacancy {
	id, err := uuid.Parse(p.Id)
	if err != nil {
		return nil
	}
	return &outPkg.Vacancy{
		ID:           id,
		Title:        p.Title,
		Description:  p.Description,
		Link:         p.Link,
		Company:      p.Company,
		Salary:       p.Salary,
		Location:     p.Location,
		Requirements: p.Requirements,
		About:        p.About,
	}
}

func mapModelToProto(m *outPkg.Vacancy) *search.Vacancy {
	return &search.Vacancy{
		Id:           m.ID.String(),
		Title:        m.Title,
		Description:  m.Description,
		Link:         m.Link,
		Company:      m.Company,
		Salary:       m.Salary,
		Location:     m.Location,
		Requirements: m.Requirements,
		About:        m.About,
	}
}

func (h *SearchHandler) SetVacancies(ctx context.Context, req *search.SetVacanciesRequest) (*search.SetVacanciesResponse, error) {
	id, err := uuid.Parse(req.Result.TaskId)
	if err != nil {
		return nil, err
	}
	vacancies := make([]outPkg.Vacancy, len(req.Result.Vacancies))
	for i, v := range req.Result.Vacancies {
		if v != nil {
			vacancies[i] = *mapProtoToModel(v)
		}
	}
	result := outPkg.ScrapeResult{
		TaskID:    id,
		Vacancies: vacancies,
	}
	err = h.usecase.SetVacancies(ctx, result)
	if err != nil {
		return nil, err
	}
	return &search.SetVacanciesResponse{}, nil
}

func (h *SearchHandler) GetVacancies(ctx context.Context, req *search.GetVacanciesRequest) (*search.GetVacanciesResponse, error) {
	if req == nil || req.Filter == nil {
		h.logger.Error("received empty request or nil filter")
		return nil, fmt.Errorf("filter is required")
	}
	filter := outPkg.VacancyFilter{
		Query:        req.Filter.Query,
		Location:     req.Filter.Location,
		Requirements: req.Filter.Requirements,
		Limit:        int(req.Filter.Limit),
		Offset:       int(req.Filter.Offset),
	}
	vacancies, err := h.usecase.GetVacancies(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]*search.Vacancy, 0, len(vacancies))
	for _, v := range vacancies {
		if v == nil {
			h.logger.Warn("found nil vacancy in search results")
			continue
		}
		result = append(result, mapModelToProto(v))
	}
	return &search.GetVacanciesResponse{
			Vacancies: result,
		},
		nil
}
