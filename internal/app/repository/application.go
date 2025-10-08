package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"lab12/internal/app/ds"
	"strings"
)

func (r *ApplicationModel) GetUserDraftManuscript(userID uint) (*ds.Manuscript, error) {
	var manuscript ds.Manuscript
	err := r.db.
		Preload("Letters.Letter").
		Where("user_id = ? AND status = ?", userID, "draft").
		Order("created_at asc").
		First(&manuscript).Error

	if err == nil {
		return &manuscript, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// если черновик не найден — создаём новый
	manuscript = ds.Manuscript{
		Status:    "draft",
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := r.db.Create(&manuscript).Error; err != nil {
		return nil, err
	}

	// и подгружаем сразу с буквами (хотя их пока нет)
	r.db.Preload("Letters.Letter").First(&manuscript, manuscript.ID)

	return &manuscript, nil
}

func (r *ApplicationModel) GetLetters() ([]ds.Letter, error) {
	var letters []ds.Letter

	err := r.db.Where("is_active = ?", true).Find(&letters).Error
	if err != nil {
		return nil, err
	}

	if len(letters) == 0 {
		return nil, fmt.Errorf("массив писем пустой")
	}

	return letters, nil
}

func (r *ApplicationModel) GetManuscript(id uint) (*ds.Manuscript, error) {
	var manuscript ds.Manuscript
	err := r.db.
		Preload("Letters.Letter").
		Where("id = ?", id).
		First(&manuscript).Error
	if err != nil {
		return nil, fmt.Errorf("рукопись не найдена")
	}

	if manuscript.Status == "deleted" {
		return nil, fmt.Errorf("рукопись удалена")
	}

	return &manuscript, nil
}

func (r *ApplicationModel) GetLetter(id uint) (ds.Letter, error) {
	var letter ds.Letter
	err := r.db.Where("id = ?", id).First(&letter).Error
	if err != nil {
		return ds.Letter{}, fmt.Errorf("письмо не найдено")
	}
	return letter, nil
}

func (r *ApplicationModel) GetLettersByName(name string) ([]ds.Letter, error) {
	var letters []ds.Letter
	err := r.db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%").Find(&letters).Error
	if err != nil {
		return nil, err
	}
	return letters, nil
}

func (r *ApplicationModel) GetManuscriptCount(manuscriptID uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.ManuscriptLetter{}).Where("manuscript_id = ?", manuscriptID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ApplicationModel) GetManuscriptTotalCount(id uint) (int, error) {
	manu, err := r.GetManuscript(id)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, item := range manu.Letters {
		total += item.Quantity
	}
	return total, nil
}

func (r *ApplicationModel) GetManuscriptWithLetters(id uint) (ds.Manuscript, error) {
	var manuscript ds.Manuscript
	err := r.db.Preload("Letters.Letter").First(&manuscript, id).Error
	return manuscript, err
}

func (r *ApplicationModel) UpdateLetterQuantity(manuscriptID, letterID uint, quantity int) error {
	return r.db.Model(&ds.ManuscriptLetter{}).
		Where("manuscript_id = ? AND letter_id = ?", manuscriptID, letterID).
		Update("quantity", quantity).Error
}

func (r *ApplicationModel) AddLetterToManuscript(manuscriptID uint, letterID uint) error {
	var manuscript ds.Manuscript
	if err := r.db.Where("id = ? AND status != 'deleted'", manuscriptID).First(&manuscript).Error; err != nil {
		return fmt.Errorf("рукопись не найдена или удалена")
	}

	var letter ds.Letter
	if err := r.db.Where("id = ? AND is_active = true", letterID).First(&letter).Error; err != nil {
		return fmt.Errorf("письмо не найдено")
	}

	var manuscriptLetter ds.ManuscriptLetter
	err := r.db.Where("manuscript_id = ? AND letter_id = ?", manuscriptID, letterID).
		First(&manuscriptLetter).Error

	if err == nil {
		// Запись уже есть
		if manuscriptLetter.IsActive {
			// просто увеличиваем количество
			manuscriptLetter.Quantity++
		} else {
			// восстанавливаем
			manuscriptLetter.IsActive = true
			manuscriptLetter.Quantity = 1
		}
		return r.db.Save(&manuscriptLetter).Error
	}

	// если записи нет — создаём
	newLetter := ds.ManuscriptLetter{
		ManuscriptID: manuscriptID,
		LetterID:     letterID,
		Quantity:     1,
		IsActive:     true,
	}
	return r.db.Create(&newLetter).Error
}

func (r *ApplicationModel) RemoveLetterFromManuscript(manuscriptID, letterID uint) error {
	res := r.db.Where("manuscript_id = ? AND letter_id = ?", manuscriptID, letterID).
		Delete(&ds.ManuscriptLetter{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("буква не найдена в рукописи")
	}
	return nil
}

func (r *ApplicationModel) GetManuscripts() ([]ds.Manuscript, error) {
	var manuscripts []ds.Manuscript

	err := r.db.Find(&manuscripts).Error
	if err != nil {
		return nil, err
	}

	if len(manuscripts) == 0 {
		return nil, fmt.Errorf("массив рукописей пустой")
	}

	return manuscripts, nil
}

func (r *ApplicationModel) SearchManuscriptsByStatus(status string) ([]ds.Manuscript, error) {
	var manuscripts []ds.Manuscript
	err := r.db.Where("status ILIKE ?", "%"+status+"%").Find(&manuscripts).Error
	if err != nil {
		return nil, err
	}
	return manuscripts, nil
}

func (r *ApplicationModel) GetManuscriptsByTitle(title string) ([]ds.Manuscript, error) {
	var manuscripts []ds.Manuscript
	err := r.db.Joins("JOIN manuscript_letters ON manuscripts.id = manuscript_letters.manuscript_id").
		Joins("JOIN letters ON manuscript_letters.letter_id = letters.id").
		Where("LOWER(letters.name) LIKE ?", "%"+strings.ToLower(title)+"%").
		Distinct().Find(&manuscripts).Error
	if err != nil {
		return nil, err
	}
	return manuscripts, nil
}

func (r *ApplicationModel) DeleteManuscript(id uint) error {
	// Выполняем SQL напрямую — обновляем статус
	result := r.db.Exec("UPDATE manuscripts SET status = 'deleted' WHERE id = ?", id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("рукопись не найдена")
	}
	return nil
}
