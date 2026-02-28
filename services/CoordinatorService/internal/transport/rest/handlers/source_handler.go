package handler

import (
	"log/slog"
	"net/http"

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

// GetAvailableSources godoc
// @Summary      Get available sources
// @Description  Returns a list of all available data sources including their IDs and names
// @Tags         sources
// @Accept       json
// @Produce      json
// @Success      200  {object}  ListSourcesResponse  "Successful response with source list"
// @Failure      500  {object}  map[string]string    "Internal server error"
// @Router       /sources [get]
func (h *Handler) GetAvailableSources(c *gin.Context) {
	ctx := c.Request.Context()

	sources, total, err := h.uc.GetAvailableSources(ctx)
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
		Sources: respSources,
		Total:   total,
	})
}
