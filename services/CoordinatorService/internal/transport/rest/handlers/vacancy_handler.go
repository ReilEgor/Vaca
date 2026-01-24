package handler

import (
	"github.com/gin-gonic/gin"
)

type VacancyResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Company string `json:"company"`
	Salary  string `json:"salary"`
	Link    string `json:"link"`
}

type SearchVacanciesResponse struct {
	Items []VacancyResponse `json:"items"`
	Total int64             `json:"total"`
}

func (h *Handler) GetVacancies(c *gin.Context) {

}
