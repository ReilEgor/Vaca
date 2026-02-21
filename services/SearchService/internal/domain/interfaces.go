package domain

import (
	"context"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	search "github.com/ReilEgor/Vaca/services/SearchService/api/proto/gen/elastic"
)

//go:generate mockery --name SearchRepository --output ../mocks/domain --outpkg domain --case=underscore
type SearchRepository interface {
	SetVacancies(ctx context.Context, vacancies outPkg.ScrapeResult) error
	GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error)
}

type SearchUsecase interface {
	SetVacancies(ctx context.Context, vacancies outPkg.ScrapeResult) error
	GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error)
}

type SearchHandler interface {
	SetVacancies(ctx context.Context, req *search.SetVacanciesRequest) (*search.SetVacanciesResponse, error)
	GetVacancies(ctx context.Context, req *search.GetVacanciesRequest) (*search.GetVacanciesResponse, error)
}
