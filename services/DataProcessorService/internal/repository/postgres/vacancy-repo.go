package postgres

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VacancyRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewVacancyRepository(db *pgxpool.Pool) domain.VacancyRepository {
	return &VacancyRepository{db: db, logger: slog.With(slog.String("component", "vacancyRepository"))}
}

const insertVacancyQuery = `
    INSERT INTO vacancies (title, company, location, salary, description, url, requirements, about, task_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT (url) DO NOTHING; 
`

func (r *VacancyRepository) SaveBatch(ctx context.Context, result outPkg.ScrapeResult) error {
	for _, vacancy := range result.Vacancies {
		_, err := r.db.Exec(ctx, insertVacancyQuery,
			vacancy.Title,
			vacancy.Company,
			vacancy.Location,
			vacancy.Salary,
			vacancy.Description,
			vacancy.Link,
			vacancy.Requirements,
			vacancy.About,
			result.TaskID,
		)
		if err != nil {
			return fmt.Errorf("save result: %w", err)
		}
	}
	return nil
}
