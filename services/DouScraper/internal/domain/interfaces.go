package domain

import (
	"context"

	outPkg "github.com/ReilEgor/Vaca/pkg"
)

//go:generate mockery --name ScraperUsecase --output ../mocks/domain --outpkg domain --case=underscore
type ScraperUsecase interface {
	Execute(ctx context.Context, task outPkg.ScrapeTask) error
}

//go:generate mockery --name TaskSubscriber --output ../mocks/domain --outpkg domain --case=underscore
type TaskSubscriber interface {
	Listen(ctx context.Context) error
}

//go:generate mockery --name ResultPublisher --output ../mocks/domain --outpkg domain --case=underscore
type ResultPublisher interface {
	PublishResults(ctx context.Context, vacancy []outPkg.Vacancy) error
}
