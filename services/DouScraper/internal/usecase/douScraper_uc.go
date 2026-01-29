package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/DouScraper/internal/broker/rabbitmq"
	"github.com/gocolly/colly"
)

type DouInteractor struct {
	logger    *slog.Logger
	publisher *rabbitmq.Publisher
}

func NewDouInteractor(publisher *rabbitmq.Publisher) *DouInteractor {
	return &DouInteractor{
		logger:    slog.With(slog.String("component", "DouInteractor")),
		publisher: publisher,
	}
}

func (i *DouInteractor) Execute(ctx context.Context, task outPkg.ScrapeTask) error {
	i.logger.Info("starting DouInteractor",
		slog.String("url", task.ID.String()),
		slog.Any("task", task))
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*jobs.dou.ua*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})
	c.OnError(func(r *colly.Response, err error) {
		i.logger.Error("Request URL failed",
			slog.String("url", r.Request.URL.String()),
			slog.Int("status", r.StatusCode),
			slog.Any("error", err))
	})

	detailCollector := c.Clone()

	var mu sync.Mutex
	var foundVacancies []outPkg.Vacancy
	c.OnHTML("li.l-vacancy", func(e *colly.HTMLElement) {
		vacancyURL := e.ChildAttr("a.vt", "href")

		ctx := colly.NewContext()
		ctx.Put("URL", vacancyURL)

		detailCollector.Request("GET", vacancyURL, nil, ctx, nil)
	})

	detailCollector.OnHTML("div.l-vacancy", func(e *colly.HTMLElement) {
		sections := e.DOM.Find("div.vacancy-section")
		i.logger.Info("Found a vacancy element!")
		vacancy := outPkg.Vacancy{
			Link:         e.Response.Ctx.Get("URL"),
			Description:  e.ChildText("div.l-t"),
			Title:        e.ChildText("h1.g-h2"),
			Company:      e.ChildText("div.l-n a"),
			City:         e.ChildText("span.place bi bi-geo-alt-fill"),
			Requirements: sections.Eq(1).Text(),
			About:        sections.Eq(2).Text(),
		}
		mu.Lock()
		foundVacancies = append(foundVacancies, vacancy)
		mu.Unlock()
	})

	err := c.Visit("https://jobs.dou.ua/vacancies/?search=" + strings.Join(task.Keyword, "+"))
	if err != nil {
		return fmt.Errorf("failed to visit: %w", err)
	}
	c.Wait()
	detailCollector.Wait()

	i.logger.Info("finished scraping", slog.Int("total", len(foundVacancies)))

	fmt.Printf("%+v\n", foundVacancies)

	result := outPkg.ScrapeResult{
		TaskID:    task.ID,
		Vacancies: foundVacancies,
	}
	err = i.publisher.PublishResults(ctx, result)
	if err != nil {
		return err
	}
	return nil
}
