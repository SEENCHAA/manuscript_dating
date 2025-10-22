package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"lab1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const (
	MinioEndpoint = "127.0.0.1:9000"
	MinioBucket   = "manuscripts"
)

func (h *ApplicationController) ApiRegisterUser(c *gin.Context) {
	var u ds.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.RegisterUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": u.ID})
}

func (h *ApplicationController) ApiLoginUser(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	user, err := h.ApplicationModel.AuthenticateUser(data.Username, data.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "username": user.Username})
}

func (h *ApplicationController) ApiLogoutUser(c *gin.Context) {
	// Логика деавторизации (например, удаление токена/сессии)
	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

func (h *ApplicationController) ApiGetUser(c *gin.Context) {
	user, err := h.ApplicationModel.GetUser(1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *ApplicationController) ApiUpdateUser(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.UpdateUser(1, data.Username, data.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// === 6-12. LETTERS (Методы) ===

func (h *ApplicationController) ApiGetLetters(c *gin.Context) {
	filter := c.Query("filter")
	letters, err := h.ApplicationModel.GetLettersFiltered(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, letters)
}

func (h *ApplicationController) ApiGetLetter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	letter, err := h.ApplicationModel.GetLetter(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "letter not found"})
		return
	}
	c.JSON(http.StatusOK, letter)
}

func (h *ApplicationController) ApiCreateLetter(c *gin.Context) {
	var body ds.Letter
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.CreateLetter(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

func (h *ApplicationController) ApiUpdateLetter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body ds.Letter
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.UpdateLetter(uint(id), &body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ApplicationController) ApiDeleteLetter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	// Здесь должна быть логика удаления изображения MinIO перед удалением записи
	if err := h.ApplicationModel.DeleteLetter(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *ApplicationController) ApiUploadImage(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Получаем ID услуги (письма)
	idStr := c.Param("id")
	letterID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid letter ID"})
		return
	}

	// 2. Парсим загруженный файл
	file, err := c.FormFile("image") // 'image' — имя поля в форме Postman
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file (use field name 'image'): " + err.Error()})
		return
	}

	// Генерируем уникальное имя файла: [ID_Услуги]_[UUID].[расширение]
	filename := fmt.Sprintf("%d_%s%s", letterID, uuid.New().String(), filepath.Ext(file.Filename))

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// 3. Загрузка файла в MinIO
	if h.MinioClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO client not initialized"})
		return
	}

	// Выполнение загрузки объекта
	_, err = h.MinioClient.PutObject(ctx, MinioBucket, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})

	if err != nil {
		// Вывод ошибки MinIO для отладки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MinIO upload failed: " + err.Error()})
		return
	}

	// 4. Формируем URL и обновляем БД
	// Формат URL: http://127.0.0.1:9000/manuscripts/filename
	imageURL := fmt.Sprintf("http://%s/%s/%s", MinioEndpoint, MinioBucket, filename)

	if err := h.ApplicationModel.UpdateLetterImageURL(uint(letterID), imageURL); err != nil {
		// В продакшене здесь должна быть логика удаления файла из MinIO (откат)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "Image successfully uploaded and link updated",
		"image_url": imageURL,
	})
}

func (h *ApplicationController) ApiAddLetterToDraft(c *gin.Context) {
	lid, _ := strconv.Atoi(c.Param("id"))
	var body struct{ Quantity int }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	// В реальном приложении здесь нужно получить UserID из сессии/токена
	userID := uint(1)
	if err := h.ApplicationModel.AddLetterToDraft(userID, uint(lid), body.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "letter added to draft", "letter_id": lid})
}

// === 13-21. MANUSCRIPTS & M-M (Методы) ===

func (h *ApplicationController) ApiGetBasketStatus(c *gin.Context) {
	// В реальном приложении здесь нужно получить UserID из сессии/токена
	userID := uint(1)
	draftID, letterCount, err := h.ApplicationModel.GetBasketStatus(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found or " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"manuscript_id": draftID, "letter_count": letterCount})
}

func (h *ApplicationController) ApiGetManuscripts(c *gin.Context) {
	status := c.Query("status")
	start := c.Query("start")
	end := c.Query("end")
	manuscripts, err := h.ApplicationModel.FilterManuscripts(status, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, manuscripts)
}

func (h *ApplicationController) ApiGetManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	manuscript, err := h.ApplicationModel.GetManuscript(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "manuscript not found"})
		return
	}
	c.JSON(http.StatusOK, manuscript)
}

func (h *ApplicationController) ApiUpdateManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var data ds.Manuscript
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.UpdateManuscript(uint(id), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ApplicationController) ApiSubmitManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.ApplicationModel.SubmitManuscript(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

func (h *ApplicationController) ApiModerateManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var data struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	var err error
	moderatorID := uint(1) // Должен быть получен из сессии/токена модератора

	switch data.Status {
	case "finished":
		// PUT завершить модератором (с расчетами)
		err = h.ApplicationModel.FinishManuscript(uint(id), moderatorID)
	case "rejected":
		// PUT отклонить модератором
		err = h.ApplicationModel.RejectManuscript(uint(id), moderatorID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. Must be 'finished' or 'rejected'"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": data.Status})
}

func (h *ApplicationController) ApiDeleteManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.ApplicationModel.DeleteManuscript(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *ApplicationController) ApiUpdateMMQuantity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	lid, _ := strconv.Atoi(c.Param("lid"))
	var body struct{ Quantity int }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.UpdateMMQuantity(uint(id), uint(lid), body.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ApplicationController) ApiDeleteMM(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	lid, _ := strconv.Atoi(c.Param("lid"))
	if err := h.ApplicationModel.DeleteMM(uint(id), uint(lid)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
