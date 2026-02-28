package handler

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc     domain.CoordinatorUsecase
	logger *slog.Logger
}

const (
	componentHandler = "handler"

	routeAPIV1            = "/api/v1"
	routeTasks            = "/tasks"
	routeTaskByID         = "/:id"
	routeVacancies        = "/vacancies"
	routeAvailableSources = "/sources"
)

func NewHandler(uc domain.CoordinatorUsecase) *Handler {
	return &Handler{
		uc:     uc,
		logger: slog.With(slog.String("component", componentHandler)),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	api := router.Group(routeAPIV1)
	{
		tasks := api.Group(routeTasks)
		{
			tasks.POST("/", h.CreateTask)
			tasks.GET(routeTaskByID, h.GetTaskStatus)
		}
		api.GET(routeVacancies, h.GetVacancies)
		api.GET(routeAvailableSources, h.GetAvailableSources)
	}
}
