package handler

import (
	"log/slog"
	"net/http"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SourceResponse struct {
	ID   uuid.UUID `json:"id" binding:"required" example:"1"`
	Name string    `json:"name" binding:"required" example:"Dou.ua"`
}

type ListSourcesResponse struct {
	Sources []SourceResponse `json:"sources" binding:"required,dive"`
	Total   int64            `json:"total" binding:"required,gte=0"`
}

func (h *Handler) GetAvailableSources(c *gin.Context) {
	sources, total, err := h.uc.GetAvailableSources(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get available sources", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToGetSources.Error()})
		return
	}

	c.JSON(http.StatusOK, ListSourcesResponse{
		Sources: mapSourcesToResponse(sources),
		Total:   total,
	})
}

func mapSourcesToResponse(sources []outPkg.Source) []SourceResponse {
	result := make([]SourceResponse, len(sources))
	for i, s := range sources {
		result[i] = SourceResponse{
			ID:   s.ID,
			Name: s.Name,
		}
	}
	return result
}
