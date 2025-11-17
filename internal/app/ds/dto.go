package ds

import (
	"time"
)

// ResponseLetterDetail - для представления письма в составе рукописи
type ResponseLetterDetail struct {
	LetterID    uint   `json:"letter_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Quantity    int    `json:"quantity"` // Количество этого письма в рукописи
}

// ManuscriptResponse - для ответа на запрос GET /api/manuscripts/{id}
type ManuscriptResponse struct {
	ID               uint      `json:"id"`
	UserID           uint      `json:"user_id"`
	Username         string    `json:"username"` // Поле для имени пользователя (загружается отдельно)
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	CalculatedPeriod string    `json:"calculated_period,omitempty"`
	ManuscriptText   string    `json:"manuscript_text"`

	// Исключаем GORM-связи, чтобы избежать рекурсии.
	Letters []ResponseLetterDetail `json:"letters"`
}

// Для отображения одной буквы в списке заявки
type ManuscriptListItemLetterResponse struct {
	ID          uint   `json:"letter_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PeriodStart int    `json:"period_start"`
	PeriodEnd   int    `json:"period_end"`
	Quantity    int    `json:"quantity"` // Важно: количество из связующей таблицы ManuscriptLetter
}

// Структура для одного элемента в списке заявок (GET /api/manuscripts)
type ManuscriptListItemResponse struct {
	ID uint `json:"id"`

	// Информация о Создателе
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`

	// Информация о Модераторе (опционально)
	ModeratorID   *uint   `json:"moderator_id,omitempty"`
	ModeratorName *string `json:"moderator_name,omitempty"` // Будет nil, если модератор не назначен

	Status           string `json:"status"`
	CalculatedPeriod string `json:"calculated_period"`

	// ПРИМЕЧАНИЕ: Если поле manuscript_text существует
	ManuscriptText string `json:"manuscript_text"`

	// Вложенные письма
	Letters []ManuscriptListItemLetterResponse `json:"letters"`
}
