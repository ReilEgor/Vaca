package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
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

type Section struct {
	Title   string
	Content []string
}

func (i *DouInteractor) Execute(ctx context.Context, task outPkg.ScrapeTask) error {
	i.logger.Info("starting DouInteractor", slog.String("url", task.ID.String()))

	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 ..."),
	)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*jobs.dou.ua*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	var (
		mu             sync.Mutex
		foundVacancies []outPkg.Vacancy
	)

	detailCollector := c.Clone()

	c.OnHTML("li.l-vacancy", func(e *colly.HTMLElement) {
		vacancyURL := e.ChildAttr("a.vt", "href")
		absURL := e.Request.AbsoluteURL(vacancyURL)

		if absURL != "" {
			ctxDetail := colly.NewContext()
			ctxDetail.Put("Company", e.ChildText("a.company"))

			detailCollector.Request("GET", absURL, nil, ctxDetail, nil)
		}
	})

	detailCollector.OnHTML("div.l-vacancy", func(e *colly.HTMLElement) {
		var requirements, about []string

		e.DOM.Find("h2, h3, p").Each(func(_ int, s *goquery.Selection) {
			text := strings.ToLower(s.Text())
			isReq := strings.Contains(text, "requirements") || strings.Contains(text, "вимоги")
			isAbout := strings.Contains(text, "about us") || strings.Contains(text, "про нас") || strings.Contains(text, "we offer")

			if isReq || isAbout {
				var target *[]string
				if isReq {
					target = &requirements
				} else {
					target = &about
				}

				next := s.Next()
				if next.Is("ul") {
					next.Find("li").Each(func(_ int, li *goquery.Selection) {
						*target = append(*target, "--"+strings.TrimSpace(li.Text()))
					})
				} else {
					*target = append(*target, strings.TrimSpace(next.Text()))
				}
			}
		})

		vacancy := outPkg.Vacancy{
			Link:         e.Request.URL.String(),
			Title:        strings.TrimSpace(e.ChildText("h1.g-h2")),
			Company:      e.Response.Ctx.Get("Company"),
			Location:     strings.TrimSpace(e.ChildText("span.place")),
			Description:  strings.TrimSpace(e.ChildText("div.b-typo.vacancy-section")),
			Requirements: strings.Join(requirements, "\n"),
			About:        strings.Join(about, "\n"),
		}

		mu.Lock()
		foundVacancies = append(foundVacancies, vacancy)
		mu.Unlock()
	})

	searchURL := fmt.Sprintf("https://jobs.dou.ua/vacancies/?search=%s", url.QueryEscape(strings.Join(task.Keyword, " ")))
	if err := c.Visit(searchURL); err != nil {
		return fmt.Errorf("failed to visit: %w", err)
	}

	c.Wait()
	detailCollector.Wait()

	return i.publisher.PublishResults(ctx, outPkg.ScrapeResult{TaskID: task.ID, Vacancies: foundVacancies})
}
