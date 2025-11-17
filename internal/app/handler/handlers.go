package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"lab31/internal/app/ds"
	"lab31/internal/app/role" // Предполагаем, что role.Role — это string/int

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

const (
	MinioEndpoint = "127.0.0.1:9000"
	MinioBucket   = "manuscripts"
	jwtPrefix     = "Bearer "
)

// ===========================================
// DTO & RESPONSE MODELS (для документации)
// ===========================================

// ErrorResponse используется для стандартных ответов с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"invalid json or record not found"`
}

// StatusResponse используется для стандартных ответов с успешным статусом
type StatusResponse struct {
	Status string `json:"status" example:"updated"`
}

// UserRegisterInput определяет тело запроса для регистрации
type UserRegisterInput struct {
	Username string `json:"username" example:"newuser"`
	Password string `json:"password" example:"securepassword"`
	// Роль не включается, так как назначается сервером
}

// RegisterResponse определяет тело ответа после регистрации
type RegisterResponse struct {
	UserID uint `json:"user_id" example:"42"`
}

// LoginRequest определяет тело запроса для входа
type LoginRequest struct {
	Username string `json:"username" example:"testuser"`
	Password string `json:"password" example:"testpassword"`
}

// LoginResponse определяет тело ответа после успешного входа
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Role  string `json:"role" example:"Buyer"`
}

// UpdateUserInput определяет тело запроса для обновления пользователя
type UpdateUserInput struct {
	Username string `json:"username" example:"updateduser"`
	Password string `json:"password" example:"newsecurepassword"`
}

// UpdateQuantityRequest определяет тело запроса для изменения количества письма
type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" example:"3"`
}

// BasketStatusResponse определяет тело ответа для статуса корзины
type BasketStatusResponse struct {
	ManuscriptID uint `json:"manuscript_id" example:"10"`
	LetterCount  int  `json:"letter_count" example:"5"`
}

// ModerateRequest определяет тело запроса для модерации
type ModerateRequest struct {
	Status string `json:"status" example:"finished"` // Должен быть 'finished' или 'rejected'
}

// ImageUploadResponse определяет тело ответа после загрузки изображения
type ImageUploadResponse struct {
	Status   string `json:"status" example:"Image successfully uploaded and link updated"`
	ImageURL string `json:"image_url" example:"http://127.0.0.1:9000/manuscripts/10_uuid.jpg"`
}

// ===========================================
// MIDDLEWARE
// ===========================================

// ===========================================
// USERS (1-5)
// ===========================================

// ApiRegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя с ролью Buyer по умолчанию.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body UserRegisterInput true "Данные нового пользователя"
// @Success 201 {object} RegisterResponse "Успешная регистрация"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} ErrorResponse "Ошибка сервера (например, пользователь уже существует)"
// @Router /users/register [post]
func (h *ApplicationController) ApiRegisterUser(c *gin.Context) {
	var u ds.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	// Роль по умолчанию - Buyer (0)
	u.Role = role.Buyer

	if err := h.ApplicationModel.RegisterUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": u.ID})
}

// ApiLoginUser godoc
// @Summary Аутентификация пользователя
// @Description Вход в систему, проверка учетных данных и выдача JWT-токена.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param login body LoginRequest true "Учетные данные пользователя"
// @Success 200 {object} LoginResponse "Успешный вход и токен"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ErrorResponse "Ошибка генерации токена"
// @Router /users/login [post]
func (h *ApplicationController) ApiLoginUser(c *gin.Context) {
	var data LoginRequest
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	user, err := h.ApplicationModel.AuthenticateUser(data.Username, data.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Генерация JWT
	userUUID := uuid.New()
	expirationTime := time.Now().Add(h.Config.JWT.ExpiresIn)
	claims := &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   user.ID,
		UserUUID: userUUID,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(h.Config.JWT.SigningMethod), claims)
	tokenString, err := token.SignedString([]byte(h.Config.JWT.Token))

	if err != nil {
		logrus.Errorf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role.String()})
}

// ApiLogoutUser godoc
// @Summary Выход из системы
// @Description Добавляет текущий JWT-токен в черный список Redis. Требуется авторизация.
// @Tags Users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешный выход"
// @Failure 401 {object} ErrorResponse "Токен не предоставлен/недействителен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера (Redis)"
// @Router /users/logout [post]
func (h *ApplicationController) ApiLogoutUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	jwtStr := strings.TrimPrefix(authHeader, jwtPrefix)
	if jwtStr == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		return
	}

	var claims ds.JWTClaims
	_, _, err := new(jwt.Parser).ParseUnverified(jwtStr, &claims) // ИСПРАВЛЕНИЕ: ПЕРЕМЕННАЯ 'token' ЗАМЕНЕНА НА '_'
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	expiresAt := time.Until(claims.ExpiresAt.Time)

	if err := h.RedisClient.AddToBlacklist(c, jwtStr, expiresAt); err != nil {
		logrus.Errorf("failed to add token to blacklist: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// ApiGetUser godoc
// @Summary Получить данные пользователя
// @Description Возвращает данные пользователя по его ID (заглушка: всегда ID 1). Требуется авторизация.
// @Tags Users
// @Produce  json
// @Param id path int true "ID пользователя"
// @Security ApiKeyAuth
// @Success 200 {object} ds.User "Данные пользователя"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Router /users/{id} [get]
func (h *ApplicationController) ApiGetUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	user, err := h.ApplicationModel.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ApiUpdateUser godoc
// @Summary Обновить данные пользователя
// @Description Обновляет имя и/или пароль пользователя (заглушка: всегда ID 1). Требуется авторизация.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "ID пользователя"
// @Param data body UpdateUserInput true "Новые данные пользователя"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное обновление"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /users/{id} [put]
func (h *ApplicationController) ApiUpdateUser(c *gin.Context) {
	var data UpdateUserInput
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if err := h.ApplicationModel.UpdateUser(1, data.Username, data.Password); err != nil { // Placeholder: ID 1
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// ===========================================
// LETTERS (6-12)
// ===========================================

// ApiGetLetters godoc
// @Summary Получить список писем (Letter)
// @Description Возвращает список всех активных писем, опционально фильтруя по тексту. Доступно всем.
// @Tags Letters
// @Produce  json
// @Param filter query string false "Фильтр для поиска по имени/описанию"
// @Success 200 {array} ds.Letter "Список писем"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /letters [get]
func (h *ApplicationController) ApiGetLetters(c *gin.Context) {
	filter := c.Query("filter")
	letters, err := h.ApplicationModel.GetLettersFiltered(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, letters)
}

// ApiGetLetter godoc
// @Summary Получить конкретное письмо
// @Description Возвращает письмо по его ID. Доступно всем.
// @Tags Letters
// @Produce  json
// @Param id path int true "ID письма"
// @Success 200 {object} ds.Letter "Детали письма"
// @Failure 404 {object} ErrorResponse "Письмо не найдено"
// @Router /letters/{id} [get]
func (h *ApplicationController) ApiGetLetter(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	letter, err := h.ApplicationModel.GetLetter(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "letter not found"})
		return
	}
	c.JSON(http.StatusOK, letter)
}

// ApiCreateLetter godoc
// @Summary Создать новое письмо
// @Description Создает новое письмо. Требуется роль Manager или Admin.
// @Tags Letters
// @Accept  json
// @Produce  json
// @Param letter body ds.Letter true "Данные нового письма"
// @Security ApiKeyAuth
// @Success 201 {object} ds.Letter "Созданное письмо"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /letters [post]
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

// ApiUpdateLetter godoc
// @Summary Обновить письмо
// @Description Обновляет данные письма по ID. Требуется роль Manager или Admin.
// @Tags Letters
// @Accept  json
// @Produce  json
// @Param id path int true "ID письма"
// @Param letter body ds.Letter true "Обновленные данные письма"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное обновление"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /letters/{id} [put]
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

// ApiDeleteLetter godoc
// @Summary Удалить письмо
// @Description Удаляет письмо по ID. Требуется роль Manager или Admin.
// @Tags Letters
// @Produce  json
// @Param id path int true "ID письма"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное удаление"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /letters/{id} [delete]
func (h *ApplicationController) ApiDeleteLetter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid letter ID"})
		return
	}

	// Вызываем функцию логического удаления из репозитория
	if err := h.ApplicationModel.DeleteLetter(uint(id)); err != nil {
		// Если письмо не найдено или другая ошибка
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем 200 OK с сообщением о логическом удалении
	c.JSON(http.StatusOK, gin.H{"status": "soft_deleted"})
}

// ApiUploadImage godoc
// @Summary Загрузить изображение для письма
// @Description Загружает изображение в MinIO и сохраняет URL в базе данных для указанного письма. Требуется роль Manager или Admin.
// @Tags Letters
// @Accept  mpfd
// @Produce  json
// @Param id path int true "ID письма"
// @Param image formData file true "Файл изображения"
// @Security ApiKeyAuth
// @Success 200 {object} ImageUploadResponse "Успешная загрузка"
// @Failure 400 {object} ErrorResponse "Неверный ID или файл не предоставлен"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера (MinIO или БД)"
// @Router /letters/{id}/image [post]
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

// ApiAddLetterToDraft godoc
// @Summary Добавить письмо в черновик рукописи
// @Description Добавляет указанное письмо в текущий черновик рукописи пользователя. Если черновика нет, он создается. Требуется роль Buyer.
// @Tags Letters
// @Accept  json
// @Produce  json
// @Param id path int true "ID письма"
// @Param data body UpdateQuantityRequest true "Количество"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{} "Письмо добавлено"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /letters/{id}/manuscript [post]
func (h *ApplicationController) ApiAddLetterToDraft(c *gin.Context) {
	lid, _ := strconv.Atoi(c.Param("id"))
	var body UpdateQuantityRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// --- ИСПРАВЛЕНИЕ: ИЗВЛЕЧЕНИЕ USERID ИЗ КОНТЕКСТА GIN ---
	userIDValue, exists := c.Get("user_id")
	if !exists {
		// Этого не должно случиться, так как маршрут защищен userAuth,
		// но это необходимая проверка безопасности.
		c.JSON(http.StatusUnauthorized, ErrorResponse{"User is not authenticated or UserID missing from context"})
		return
	}

	// Приводим интерфейс к нужному типу (uint).
	// Важно, чтобы ваш middleware сохранял его как uint.
	currentUserID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{"Invalid UserID format in context"})
		return
	}
	// -------------------------------------------------------------

	// Используем ПРАВИЛЬНЫЙ ID
	if err := h.ApplicationModel.AddLetterToDraft(currentUserID, uint(lid), body.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "letter added to draft", "letter_id": lid})
}

// ===========================================
// MANUSCRIPTS & M-M (13-21)
// ===========================================

// ApiGetBasketStatus godoc
// @Summary Статус корзины/черновика
// @Description Возвращает ID текущего черновика (Manuscript) пользователя и количество писем в нем. Требуется роль Buyer.
// @Tags Manuscripts
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} BasketStatusResponse "Статус черновика"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 404 {object} ErrorResponse "Черновик не найден"
// @Router /manuscripts/basket [get]
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

// ApiGetManuscripts godoc
// @Summary Фильтрация списка рукописей
// @Description Возвращает список рукописей, отфильтрованный по статусу и/или датам. Администраторы и модераторы видят все заявки, покупатели – только свои.
// @Tags Manuscripts
// @Produce  json
// @Param status query string false "Фильтр по статусу (например, 'submitted', 'finished')"
// @Param start query string false "Начальная дата (YYYY-MM-DD)"
// @Param end query string false "Конечная дата (YYYY-MM-DD)"
// @Security ApiKeyAuth
// @Success 200 {array} ds.ManuscriptListItemResponse "Список рукописей"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /manuscripts [get]
func (h *ApplicationController) ApiGetManuscripts(c *gin.Context) {

	userIDValue, idExists := c.Get("user_id")
	roleValue, roleExists := c.Get("user_role")

	if !idExists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication context missing"})
		return
	}

	currentUserID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UserID format in context"})
		return
	}

	currentRole, ok := roleValue.(role.Role)
	if !ok {
		currentRole = role.Buyer
	}

	// 2. ОПРЕДЕЛЕНИЕ ФИЛЬТРАЦИИ ПО ID
	var filterUserID *uint // Используем указатель для опционального фильтра

	// Если пользователь - Buyer (0), применяем фильтр по его ID.
	// Если Manager (1) или Admin (2), filterUserID останется nil, и фильтрация не будет применена.
	if currentRole == role.Buyer {
		filterUserID = &currentUserID // Устанавливаем фильтр по его ID
	}

	// 3. Получение параметров запроса
	status := c.Query("status")
	start := c.Query("start")
	end := c.Query("end")

	// 4. Вызов модели с дополнительным параметром filterUserID
	// Модель должна возвращать полные структуры ds.Manuscript с Preload!
	manuscripts, err := h.ApplicationModel.FilterManuscripts(status, start, end, filterUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. МАППИНГ РЕЗУЛЬТАТОВ НА DTO
	responseList := make([]ds.ManuscriptListItemResponse, 0, len(manuscripts))

	for _, m := range manuscripts {
		// 5.1 Маппинг вложенных писем (Letters)
		letterList := make([]ds.ManuscriptListItemLetterResponse, 0, len(m.Letters))
		for _, ml := range m.Letters {
			// Проверка, что Letter (признак) был загружен через Preload
			if ml.Letter.ID == 0 {
				continue
			}

			letterList = append(letterList, ds.ManuscriptListItemLetterResponse{
				ID:          ml.Letter.ID,
				Name:        ml.Letter.Name,
				Description: ml.Letter.Description,
				PeriodStart: ml.Letter.PeriodStart,
				PeriodEnd:   ml.Letter.PeriodEnd,
				Quantity:    ml.Quantity, // ml - это ManuscriptLetter, содержит Quantity
			})
		}

		// 5.2 Определение имени модератора (обработка *uint)
		var moderatorName *string
		if m.ModeratorID != nil && m.Moderator.Username != "" {
			name := m.Moderator.Username
			moderatorName = &name
		}

		// 5.3 Заполнение основной DTO
		responseList = append(responseList, ds.ManuscriptListItemResponse{
			ID:               m.ID,
			UserID:           m.UserID,
			Username:         m.User.Username, // Имя создателя
			ModeratorID:      m.ModeratorID,
			ModeratorName:    moderatorName,
			Status:           m.Status,
			CalculatedPeriod: m.CalculatedPeriod,
			ManuscriptText:   m.ManuscriptText, // Поле из ds.Manuscript
			Letters:          letterList,
		})
	}

	// 6. Отправка ответа
	c.JSON(http.StatusOK, responseList)
}

// ApiGetManuscript godoc
// @Summary Получить конкретную рукопись
// @Description Возвращает детали рукописи по ID. Требуется авторизация (Owner, Manager, Admin).
// @Tags Manuscripts
// @Produce  json
// @Param id path int true "ID рукописи"
// @Security ApiKeyAuth
// @Success 200 {object} ds.Manuscript "Детали рукописи"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 404 {object} ErrorResponse "Рукопись не найдена"
// @Router /manuscripts/{id} [get]
func (h *ApplicationController) ApiGetManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	manuscript, err := h.ApplicationModel.GetManuscript(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "manuscript not found"})
		return
	}
	c.JSON(http.StatusOK, manuscript)
}

// ApiUpdateManuscript godoc
// @Summary Обновить рукопись
// @Description Обновляет поля рукописи по ID. Доступно только владельцу (Owner) и только в статусе 'draft'.
// @Tags Manuscripts
// @Accept  json
// @Produce  json
// @Param id path int true "ID рукописи"
// @Param data body ds.Manuscript true "Обновленные данные рукописи"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное обновление"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /manuscripts/{id} [put]
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

// ApiSubmitManuscript godoc
// @Summary Отправить рукопись на модерацию
// @Description Меняет статус рукописи с 'draft' на 'submitted'. Доступно только владельцу (Owner).
// @Tags Manuscripts
// @Produce  json
// @Param id path int true "ID рукописи"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешно отправлено"
// @Failure 400 {object} ErrorResponse "Неверный статус или ошибка БД"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Router /manuscripts/{id}/submit [put]
func (h *ApplicationController) ApiSubmitManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.ApplicationModel.SubmitManuscript(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

// ApiModerateManuscript godoc
// @Summary Модерация рукописи
// @Description Завершает ('finished') или отклоняет ('rejected') рукопись. Требуется роль Manager или Admin.
// @Tags Manuscripts
// @Accept  json
// @Produce  json
// @Param id path int true "ID рукописи"
// @Param data body ModerateRequest true "Статус модерации: 'finished' или 'rejected'"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешная модерация"
// @Failure 400 {object} ErrorResponse "Неверный статус или ошибка БД"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Router /manuscripts/{id}/moderation [put]
func (h *ApplicationController) ApiModerateManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var data ModerateRequest
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

// ApiDeleteManuscript godoc
// @Summary Удалить рукопись
// @Description Удаляет рукопись по ID. Доступно владельцу (Owner) в статусе 'draft' или Admin.
// @Tags Manuscripts
// @Produce  json
// @Param id path int true "ID рукописи"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное удаление"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /manuscripts/{id} [delete]
func (h *ApplicationController) ApiDeleteManuscript(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.ApplicationModel.DeleteManuscript(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ApiUpdateMMQuantity godoc
// @Summary Обновить количество письма в рукописи
// @Description Обновляет количество конкретного письма (Letter) внутри рукописи (Manuscript). Доступно владельцу (Owner) в статусе 'draft'.
// @Tags Manuscripts
// @Accept  json
// @Produce  json
// @Param id path int true "ID рукописи"
// @Param lid path int true "ID письма"
// @Param data body UpdateQuantityRequest true "Новое количество"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное обновление"
// @Failure 400 {object} ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /manuscripts/{id}/letters/{lid} [put]
func (h *ApplicationController) ApiUpdateMMQuantity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	lid, _ := strconv.Atoi(c.Param("lid"))
	var body UpdateQuantityRequest
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

// ApiDeleteMM godoc
// @Summary Удалить письмо из рукописи
// @Description Удаляет связь между письмом (Letter) и рукописью (Manuscript). Доступно владельцу (Owner) в статусе 'draft'.
// @Tags Manuscripts
// @Produce  json
// @Param id path int true "ID рукописи"
// @Param lid path int true "ID письма"
// @Security ApiKeyAuth
// @Success 200 {object} StatusResponse "Успешное удаление"
// @Failure 401 {object} ErrorResponse "Неавторизован"
// @Failure 403 {object} ErrorResponse "Нет доступа"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /manuscripts/{id}/letters/{lid} [delete]
func (h *ApplicationController) ApiDeleteMM(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	lid, _ := strconv.Atoi(c.Param("lid"))
	if err := h.ApplicationModel.DeleteMM(uint(id), uint(lid)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
