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
	ID           uuid.UUID `json:"id" binding:"required"`
	Title        string    `json:"title" binding:"required,min=3"`
	Company      string    `json:"company" binding:"required"`
	Salary       string    `json:"salary,omitempty"`
	Location     string    `json:"location,omitempty"`
	Description  string    `json:"description,omitempty"`
	Link         string    `json:"link" binding:"required,url"`
	About        string    `json:"about" binding:"required"`
	Requirements string    `json:"requirements" binding:"required"`
}

type SearchVacanciesResponse struct {
	Items []VacancyResponse `json:"items" binding:"required,dive"`
	Total int64             `json:"total" binding:"required,gte=0"`
}

func (h *Handler) GetVacancies(c *gin.Context) {
	var filter outPkg.VacancyFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn("invalid vacancies filter", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidRequestBody.Error()})
		return
	}

	vacancies, total, err := h.uc.GetVacancies(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("failed to get vacancies", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToGetVacancies.Error()})
		return
	}
	c.JSON(http.StatusOK, SearchVacanciesResponse{
		Items: mapVacanciesToResponse(vacancies),
		Total: total,
	})
}

func mapVacanciesToResponse(vacancies []*outPkg.Vacancy) []VacancyResponse {
	result := make([]VacancyResponse, len(vacancies))
	for i, v := range vacancies {
		result[i] = VacancyResponse{
			ID:           v.ID,
			Title:        v.Title,
			Company:      v.Company,
			Salary:       v.Salary,
			Location:     v.Location,
			Description:  v.Description,
			Link:         v.Link,
			About:        v.About,
			Requirements: v.Requirements,
		}
	}
	return result
}
