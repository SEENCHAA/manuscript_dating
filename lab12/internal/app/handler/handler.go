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

	// ⚡ Важно: теперь нужно указать ID рукописи
	total := h.repo.GetTotalManuscriptCount(1) // пока фиксированно 1

	data := gin.H{
		"Signs":      signsToShow,
		"TotalCount": total,
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "некорректный id")
		return
	}

	manuscript, err := h.repo.GetManuscript(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	data := gin.H{
		"Signs":      manuscript.Signs,
		"TotalCount": h.repo.GetTotalManuscriptCount(id),
	}

	c.HTML(http.StatusOK, "manuscript.html", data)
}
