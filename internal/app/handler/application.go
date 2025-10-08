package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"lab12/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *ApplicationController) RegisterController(router *gin.Engine) {
	// Главная страница (признаки)
	router.GET("/signs", h.GetLetters)
	router.GET("/sign/:id", h.GetLetter)

	// Рукопись (черновик)
	router.GET("/manuscript/:id", h.GetManuscriptByID)
	router.POST("/manuscript/add", h.AddLetterToManuscript)
	router.POST("/manuscript/remove/:id", h.RemoveLetterFromManuscript)
	router.POST("/manuscript/update/:id", h.UpdateLetterQuantity)

	// Удаление рукописи (логическое)
	router.POST("/manuscript/delete/:id", h.DeleteManuscript)

	// Список всех рукописей
	router.GET("/manuscripts", h.GetManuscripts)
}

func (h *ApplicationController) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *ApplicationController) errorController(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

// ---------- Признаки (буквы) ----------

func (h *ApplicationController) GetLetters(ctx *gin.Context) {
	var letters []ds.Letter
	var err error

	record := ctx.Query("record")

	if record == "" {
		letters, err = h.ApplicationModel.GetLetters()
	} else {
		letters, err = h.ApplicationModel.GetLettersByName(record)
	}

	if err != nil {
		logrus.Error(err)
	}

	// Получаем текущий черновик пользователя
	draft, err := h.ApplicationModel.GetUserDraftManuscript(1)
	if err != nil {
		logrus.Error(err)
	}

	// Считаем количество букв в черновике
	totalCount, err := h.ApplicationModel.GetManuscriptTotalCount(draft.ID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "signs.html", gin.H{
		"Signs":        letters,
		"TotalCount":   totalCount,
		"Records":      record,
		"ManuscriptID": draft.ID,
	})
}

func (h *ApplicationController) GetLetter(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверный ID"))
		return
	}

	letter, err := h.ApplicationModel.GetLetter(uint(id))
	if err != nil {
		h.errorController(ctx, http.StatusNotFound, fmt.Errorf("письмо не найдено"))
		return
	}

	ctx.HTML(http.StatusOK, "sign.html", gin.H{
		"Sign": letter,
	})
}

// ---------- Рукописи ----------

func (h *ApplicationController) GetManuscriptByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверный ID рукописи"))
		return
	}

	var manuscript *ds.Manuscript

	if id == 1 {
		// всегда загружаем черновик пользователя
		manuscript, err = h.ApplicationModel.GetUserDraftManuscript(1)
	} else {
		manuscript, err = h.ApplicationModel.GetManuscript(uint(id))
	}

	if err != nil {
		h.errorController(ctx, http.StatusNotFound, err)
		return
	}

	totalCount, err := h.ApplicationModel.GetManuscriptTotalCount(manuscript.ID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "manuscript.html", gin.H{
		"Manuscript":   manuscript,
		"ManuscriptID": manuscript.ID,
		"TotalCount":   totalCount,
	})
}

func (h *ApplicationController) AddLetterToManuscript(ctx *gin.Context) {
	letterIDStr := ctx.PostForm("id")
	letterID, err := strconv.Atoi(letterIDStr)
	if err != nil {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверный ID признака"))
		return
	}

	// получаем текущий черновик (создаст новый, если удалён)
	draft, err := h.ApplicationModel.GetUserDraftManuscript(1)
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.ApplicationModel.AddLetterToManuscript(draft.ID, uint(letterID))
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/signs")
}

func (h *ApplicationController) RemoveLetterFromManuscript(ctx *gin.Context) {
	letterIDStr := ctx.Param("id")
	letterID, err := strconv.Atoi(letterIDStr)
	if err != nil {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверный ID буквы"))
		return
	}

	// всегда берём актуальный черновик
	draft, err := h.ApplicationModel.GetUserDraftManuscript(1)
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.ApplicationModel.RemoveLetterFromManuscript(draft.ID, uint(letterID))
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/manuscript/%d", draft.ID))
}

func (h *ApplicationController) UpdateLetterQuantity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	letterID, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверный ID буквы"))
		return
	}

	quantityStr := ctx.PostForm("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверное количество"))
		return
	}

	draftManuscript, err := h.ApplicationModel.GetUserDraftManuscript(1)
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.ApplicationModel.UpdateLetterQuantity(draftManuscript.ID, uint(letterID), quantity)
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/manuscript/%d", draftManuscript.ID))
}

func (h *ApplicationController) DeleteManuscript(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorController(ctx, http.StatusBadRequest, fmt.Errorf("неверный ID рукописи"))
		return
	}

	// Вызов SQL-удаления
	err = h.ApplicationModel.DeleteManuscript(uint(id))
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/signs")
}

func (h *ApplicationController) GetManuscripts(ctx *gin.Context) {
	var manuscripts []ds.Manuscript
	var err error

	searchQuery := ctx.Query("search")
	if searchQuery == "" {
		manuscripts, err = h.ApplicationModel.GetManuscripts()
	} else {
		manuscripts, err = h.ApplicationModel.GetManuscriptsByTitle(searchQuery)
	}
	if err != nil {
		h.errorController(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "signs.html", gin.H{
		"Signs": manuscripts,
		"Query": searchQuery,
	})
}
