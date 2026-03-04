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

// GetVacancies godoc
// @Summary      Search and filter vacancies
// @Description  Get a paginated list of job vacancies filtered by title, source, or keywords
// @Tags         vacancies
// @Produce      json
// @Param        title     query     string  false  "Search by job title"
// @Param        source    query     string  false  "Filter by job source (e.g., djinni, dou)"
// @Param        limit     query     int     false  "Number of items to return" default(10)
// @Param        offset    query     int     false  "Number of items to skip" default(0)
// @Success      200       {object}  SearchVacanciesResponse
// @Failure      400       {object}  map[string]string "Invalid query parameters"
// @Failure      500       {object}  map[string]string "Internal server error"
// @Router       /vacancies [get]
func (h *Handler) GetVacancies(c *gin.Context) {
	var filter outPkg.VacancyFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn("invalid request body", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidRequestBody.Error()})
		return
	}
	ctx := c.Request.Context()

	vacancies, quantity, err := h.uc.GetVacancies(ctx, filter)
	if err != nil {
		h.logger.Error("failed to get vacancies", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrFailedToGetVacancies.Error()})
		return
	}
	respItems := make([]VacancyResponse, len(vacancies))
	for i, v := range vacancies {
		respItems[i] = VacancyResponse{
			ID:           v.ID,
			Title:        v.Title,
			Company:      v.Company,
			Salary:       v.Salary,
			Link:         v.Link,
			Location:     v.Location,
			Description:  v.Description,
			About:        v.About,
			Requirements: v.Requirements,
		}
	}

	c.JSON(http.StatusOK, SearchVacanciesResponse{
		Items: respItems,
		Total: quantity,
	})
}
