package repository

import (
	"errors"
	"time"

	"lab31/internal/app/ds"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ApplicationModel struct {
	db *gorm.DB
}

func NewApplicationModel(dsn string) (*ApplicationModel, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&ds.Letter{}, &ds.Manuscript{}, &ds.ManuscriptLetter{}, &ds.User{})
	if err != nil {
		return nil, err
	}

	return &ApplicationModel{db: db}, nil
}

// ===== LETTERS (приведены для полноты) =====

func (m *ApplicationModel) GetLettersFiltered(filter string) ([]ds.Letter, error) {
	var letters []ds.Letter
	query := m.db.Where("is_active = ?", true)
	if filter != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR details LIKE ?", "%"+filter+"%", "%"+filter+"%", "%"+filter+"%")
	}
	return letters, query.Find(&letters).Error
}

func (m *ApplicationModel) GetLetter(id uint) (*ds.Letter, error) {
	var letter ds.Letter
	return &letter, m.db.First(&letter, id).Error
}

func (m *ApplicationModel) CreateLetter(letter *ds.Letter) error {
	return m.db.Create(letter).Error
}

func (m *ApplicationModel) UpdateLetter(id uint, data *ds.Letter) error {
	return m.db.Model(&ds.Letter{}).Where("id = ?", id).Updates(data).Error
}

func (m *ApplicationModel) DeleteLetter(id uint) error {
	return m.db.Model(&ds.Letter{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// ===== MANUSCRIPTS (Исправленная логика) =====

// GetBasketStatus: Находит черновик пользователя и считает услуги.
func (m *ApplicationModel) GetBasketStatus(userID uint) (draftID uint, letterCount int, err error) {
	var ms ds.Manuscript
	// Найти черновик пользователя
	err = m.db.Where("user_id = ? AND status = ?", userID, "draft").Preload("Letters").First(&ms).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, errors.New("draft manuscript not found")
	} else if err != nil {
		return 0, 0, err
	}

	// Посчитать количество услуг (письмо-рукопись)
	return ms.ID, len(ms.Letters), nil
}

func (m *ApplicationModel) FilterManuscripts(status, start, end string, filterUserID *uint) ([]ds.Manuscript, error) {
	var manuscripts []ds.Manuscript
	query := m.db.Model(&ds.Manuscript{}).
		Preload("Letters.Letter").
		Preload("User").
		Preload("Moderator").
		Order("created_at DESC")

	// 1. УСЛОВНАЯ ФИЛЬТРАЦИЯ ПО USERID
	if filterUserID != nil {
		// Если filterUserID не nil (т.е. это Buyer), добавляем WHERE user_id
		query = query.Where("user_id = ?", *filterUserID)
	}

	// 2. Фильтрация по статусу
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 3. Фильтрация по времени (простейший вариант)
	if start != "" {
		query = query.Where("created_at >= ?", start)
	}
	if end != "" {
		query = query.Where("created_at <= ?", end)
	}

	err := query.Find(&manuscripts).Error
	if err != nil {
		return nil, err
	}

	return manuscripts, nil
}

func (m *ApplicationModel) GetManuscript(mid uint) (*ds.ManuscriptResponse, error) {
	// 1. Получаем основную запись Manuscript и имя пользователя (Username)
	var manuscriptData struct {
		ds.Manuscript
		Username string `json:"username"` // Дополнительное поле для имени
	}

	// Используем JOIN для получения Username и SELECT для получения всех полей Manuscript
	err := m.db.Model(&ds.Manuscript{}).
		Select("manuscripts.*, users.username").
		Joins("JOIN users ON users.id = manuscripts.user_id").
		First(&manuscriptData, mid).Error

	if err != nil {
		return nil, err
	}

	// 2. Получаем связанные письма (Letter) и их количество (Quantity)
	var lettersDetails []ds.ResponseLetterDetail

	// Используем JOIN с ManuscriptLetter и Letter и заполняем DTO-структуру
	err = m.db.Table("manuscript_letters").
		Select("manuscript_letters.letter_id, manuscript_letters.quantity, letters.name, letters.description, letters.image_url").
		Joins("JOIN letters ON letters.id = manuscript_letters.letter_id").
		Where("manuscript_letters.manuscript_id = ?", mid).
		Where("letters.is_active = ?", true). // Исключаем логически удаленные письма
		Scan(&lettersDetails).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 3. Собираем финальный DTO-ответ
	response := &ds.ManuscriptResponse{
		ID:               manuscriptData.ID,
		UserID:           manuscriptData.UserID,
		Username:         manuscriptData.Username,
		Status:           manuscriptData.Status,
		CreatedAt:        manuscriptData.CreatedAt,
		CalculatedPeriod: manuscriptData.CalculatedPeriod,
		ManuscriptText:   manuscriptData.ManuscriptText,
		Letters:          lettersDetails,
	}

	return response, nil
}

func (m *ApplicationModel) UpdateManuscript(id uint, data *ds.Manuscript) error {
	return m.db.Model(&ds.Manuscript{}).Where("id = ?", id).Updates(data).Error
}

func (m *ApplicationModel) SubmitManuscript(id uint) error {
	now := time.Now()
	// Проверка на обязательные поля должна быть реализована здесь (или в сервисе)
	return m.db.Model(&ds.Manuscript{}).Where("id = ? AND status = 'draft'", id).
		Updates(map[string]interface{}{"status": "submitted", "submitted_at": &now}).Error
}

// FinishManuscript: Завершение модератором (с расчетом Lab-2)
func (m *ApplicationModel) FinishManuscript(manuscriptID uint, moderatorID uint) error {
	now := time.Now()
	// Здесь должна быть реализована формула расчета (например, CalculatedPeriod)
	return m.db.Model(&ds.Manuscript{}).Where("id = ? AND status = 'submitted'", manuscriptID).Updates(map[string]interface{}{
		"status":            "finished",
		"moderator_id":      moderatorID,
		"finished_at":       &now,
		"calculated_period": "calculated_value", // Placeholder for actual calculation
	}).Error
}

// RejectManuscript: Отклонение модератором
func (m *ApplicationModel) RejectManuscript(manuscriptID uint, moderatorID uint) error {
	now := time.Now()
	return m.db.Model(&ds.Manuscript{}).Where("id = ? AND status = 'submitted'", manuscriptID).Updates(map[string]interface{}{
		"status":            "rejected",
		"moderator_id":      moderatorID,
		"finished_at":       &now,
		"calculated_period": nil,
	}).Error
}

func (m *ApplicationModel) DeleteManuscript(id uint) error {
	// Удаление (смена статуса на 'deleted')
	return m.db.Model(&ds.Manuscript{}).Where("id = ?", id).
		Update("status", "deleted").Error
}

// ===== M-M (приведены для полноты) =====

// AddLetterToDraft: Находит или создает черновик и добавляет услугу
func (m *ApplicationModel) AddLetterToDraft(userID uint, letterID uint, quantity int) error {
	var ms ds.Manuscript

	// 1. Найти или создать черновик
	err := m.db.Where("user_id = ? AND status = ?", userID, "draft").First(&ms).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Создать новый черновик
		ms = ds.Manuscript{UserID: userID, Status: "draft", CreatedAt: time.Now()}
		if err := m.db.Create(&ms).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 2. Добавить/обновить запись в M-M таблице
	var ml ds.ManuscriptLetter
	// Проверить, существует ли уже M-M запись
	if m.db.Where("manuscript_id = ? AND letter_id = ?", ms.ID, letterID).First(&ml).Error == nil {
		// Обновить существующую
		return m.db.Model(&ml).Update("quantity", quantity).Error
	}

	// Создать новую запись
	ml = ds.ManuscriptLetter{
		ManuscriptID: ms.ID,
		LetterID:     letterID,
		Quantity:     quantity,
	}
	return m.db.Create(&ml).Error
}

func (m *ApplicationModel) UpdateMMQuantity(mid, lid uint, quantity int) error {
	// Должен также обновлять порядок/значение, если они есть в ds.ManuscriptLetter
	return m.db.Model(&ds.ManuscriptLetter{}).Where("manuscript_id = ? AND letter_id = ? AND quantity > 0", mid, lid).
		Update("quantity", quantity).Error
}

func (m *ApplicationModel) DeleteMM(mid, lid uint) error {
	return m.db.Where("manuscript_id = ? AND letter_id = ?", mid, lid).Delete(&ds.ManuscriptLetter{}).Error
}

// ===== USERS (приведены для полноты) =====

func (m *ApplicationModel) RegisterUser(user *ds.User) error {
	// ХЭШИРОВАНИЕ ПАРОЛЯ перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	// Role уже должен быть установлен (по умолчанию 0 - Buyer)
	return m.db.Create(user).Error
}

func (m *ApplicationModel) AuthenticateUser(username, password string) (*ds.User, error) {
	var user ds.User

	// 1. Найти пользователя по имени
	err := m.db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Сравнить хэш пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (m *ApplicationModel) UpdateLetterImageURL(id uint, url string) error {
	// В структуре ds.Letter должно быть поле ImageURL string
	return m.db.Model(&ds.Letter{}).Where("id = ?", id).
		Update("image_url", url).Error
}

func (m *ApplicationModel) GetUser(id uint) (*ds.User, error) {
	var user ds.User
	err := m.db.First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	return &user, err
}

func (m *ApplicationModel) UpdateUser(id uint, username, password string) error {
	return m.db.Model(&ds.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"username": username,
		"password": password,
	}).Error
}
