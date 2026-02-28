package handler

import (
	"log/slog"

	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	uc     domain.CoordinatorUsecase
	logger *slog.Logger
}

func NewHandler(uc domain.CoordinatorUsecase) *Handler {
	return &Handler{
		uc:     uc,
		logger: slog.With(slog.String("component", "handler")),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("/", h.CreateTask)
			tasks.GET("/:id", h.GetTaskStatus)
		}

		api.GET("/vacancies", h.GetVacancies)
		api.GET("/sources", h.GetAvailableSources)

	}
}
