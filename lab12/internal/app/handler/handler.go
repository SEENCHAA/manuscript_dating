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

func (h *Handler) GetFeatures(c *gin.Context) {
	var featuresToShow []repository.Feature

	// Обработка поиска
	query := c.Query("search")
	if query != "" {
		featuresToShow, _ = h.repo.SearchFeatures(query)
	} else {
		featuresToShow, _ = h.repo.GetFeatures()
	}

	// Обработка добавления в заказ (только если POST)
	if c.Request.Method == "POST" {
		idStr := c.PostForm("id")
		if idStr != "" {
			id, _ := strconv.Atoi(idStr)
			h.repo.AddToOrder(id)
		}
	}

	c.HTML(http.StatusOK, "index.html", featuresToShow)
}

func (h *Handler) GetFeature(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	selected, _ := h.repo.GetFeature(id)

	// Добавление в заказ
	if c.Request.Method == "POST" {
		h.repo.AddToOrder(id)
		c.Redirect(http.StatusSeeOther, "/order")
		return
	}

	c.HTML(http.StatusOK, "detail.html", selected)
}

func (h *Handler) Order(c *gin.Context) {
	order := h.repo.GetOrder()
	c.HTML(http.StatusOK, "order.html", order)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)

	h.repo.RemoveFromOrder(id)

	c.Redirect(http.StatusSeeOther, "/order")
}
