package handler

import (
	"log/slog"
	"net/http"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VacancyResponse struct {
	ID      uuid.UUID `json:"id" binding:"required"`
	Title   string    `json:"title" binding:"required,min=3"`
	Company string    `json:"company" binding:"required"`
	Salary  string    `json:"salary,omitempty"`
	Link    string    `json:"link" binding:"required,url"`
}

type SearchVacanciesResponse struct {
	Items []VacancyResponse `json:"items" binding:"required,dive"`
	Total int64             `json:"total" binding:"required,gte=0"`
}

func (h *Handler) GetVacancies(c *gin.Context) {
	var req outPkg.VacancyFilter

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}
	ctx := c.Request.Context()

	vacancies, quantity, err := h.uc.GetVacancies(ctx, req)
	if err != nil {
		h.logger.Error("failed to get vacancies", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToGetVacancies.Error()})
		return
	}
	respItems := make([]VacancyResponse, len(vacancies))
	for i, v := range vacancies {
		respItems[i] = VacancyResponse{
			ID:      v.ID,
			Title:   v.Title,
			Company: v.Company,
			Salary:  v.Salary,
			Link:    v.Link,
		}
	}

	c.JSON(http.StatusOK, SearchVacanciesResponse{
		Items: respItems,
		Total: quantity,
	})
}
