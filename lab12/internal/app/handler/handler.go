package handler

import (
	"net/http"
	"strconv"

	"lab12/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// ====== Список признаков ======
func (h *Handler) GetSigns(c *gin.Context) {
	var signsToShow []repository.Feature

	// Поиск
	query := c.Query("search")
	if query != "" {
		signsToShow, _ = h.repo.SearchFeatures(query)
	} else {
		signsToShow, _ = h.repo.GetFeatures()
	}

	data := gin.H{
		"Signs":      signsToShow,
		"TotalCount": h.repo.GetTotalManuscriptCount(),
	}

	c.HTML(http.StatusOK, "signs.html", data)
}

// ====== Один признак ======
func (h *Handler) GetSign(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	selected, _ := h.repo.GetFeature(id)

	c.HTML(http.StatusOK, "sign.html", selected)
}

// ====== Рукопись ======
func (h *Handler) GetManuscript(c *gin.Context) {
	manuscript := h.repo.GetManuscript()
	data := gin.H{
		"Signs":      manuscript,
		"TotalCount": h.repo.GetTotalManuscriptCount(),
	}
	c.HTML(http.StatusOK, "manuscript.html", data)
}
