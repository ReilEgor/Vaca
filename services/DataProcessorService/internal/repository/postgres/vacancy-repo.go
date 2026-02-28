package postgres

import (
	"context"
	"fmt"
	"log/slog"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DataProcessorService/internal/domain"
	"github.com/jackc/pgx/v5"
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
	batch := &pgx.Batch{}
	for _, v := range result.Vacancies {
		batch.Queue(insertVacancyQuery,
			v.Title, v.Company, v.Location, v.Salary,
			v.Description, v.Link, v.Requirements, v.About,
			result.TaskID,
		)
	}

	results := r.db.SendBatch(ctx, batch)
	defer results.Close()

	for range result.Vacancies {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("save batch: %w", err)
		}
	}

	return nil
}
