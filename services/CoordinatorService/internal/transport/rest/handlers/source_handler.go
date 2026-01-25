package handler

import (
	"log/slog"
	"net/http"

	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SourceResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ListSourcesResponse struct {
	Sources  []SourceResponse `json:"sources"`
	Quantity int64            `json:"quantity"`
}

func (h *Handler) GetAvailableSources(c *gin.Context) {
	ctx := c.Request.Context()

	sources, quantity, err := h.uc.GetAvailableSources(ctx)
	if err != nil {
		h.logger.Error("failed to get available sources", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToGetSources.Error()})
		return
	}

	respSources := make([]SourceResponse, len(sources))
	for i, s := range sources {
		respSources[i] = SourceResponse{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	c.JSON(http.StatusOK, ListSourcesResponse{
		Sources:  respSources,
		Quantity: quantity,
	})
}
