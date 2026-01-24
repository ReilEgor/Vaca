package handler

import (
	"github.com/gin-gonic/gin"
)

type SourceResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	IsActive bool   `json:"is_active"`
}

type ListSourcesResponse struct {
	Sources []SourceResponse `json:"sources"`
}

func (h *Handler) GetAvailableSources(c *gin.Context) {

}
