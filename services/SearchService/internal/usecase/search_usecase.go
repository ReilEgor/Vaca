package usecase

import (
	"context"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/SearchService/internal/domain"
)

type SearchInteractor struct {
	logger *slog.Logger
	repo   domain.SearchRepository
}

func NewSearchUsecase(repo domain.SearchRepository) *SearchInteractor {
	return &SearchInteractor{
		logger: slog.With(slog.String("component", "search_uc")),
		repo:   repo,
	}
}

func (u *SearchInteractor) SetVacancies(ctx context.Context, vacancies outPkg.ScrapeResult) error {
	u.logger.Debug("SetVacancies")
	err := u.repo.SetVacancies(ctx, vacancies)
	if err != nil {
		return err
	}
	return nil
}

func (u *SearchInteractor) GetVacancies(ctx context.Context, filter outPkg.VacancyFilter) ([]*outPkg.Vacancy, error) {
	u.logger.Debug("GetVacancies")
	vacancies, err := u.repo.GetVacancies(ctx, filter)
	if err != nil {
		return nil, err
	}
	return vacancies, nil
}
