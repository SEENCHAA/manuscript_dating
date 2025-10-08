package ds

import "time"

// Пользователи
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Username    string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	IsModerator bool   `gorm:"default:false"`
}

// Буквы (letters)
type Letter struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	PeriodStart int    `gorm:"not null"`
	PeriodEnd   int    `gorm:"not null"`
	Details     string `gorm:"not null"`
	ImageURL    string
	IsActive    bool `gorm:"default:true"`
}

// Рукописи (manuscripts)
type Manuscript struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null;default:'draft'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	SubmittedAt *time.Time
	FinishedAt  *time.Time
	ModeratorID *uint

	// связи
	User    User               `gorm:"foreignKey:UserID"`
	Letters []ManuscriptLetter `gorm:"foreignKey:ManuscriptID"`
}

// Связь "Рукопись - Буква"
type ManuscriptLetter struct {
	ManuscriptID uint `gorm:"primaryKey"`
	LetterID     uint `gorm:"primaryKey"`
	Quantity     int  `gorm:"default:1"`
	IsActive     bool `gorm:"default:true"` // логическое удаление

	// связи
	Letter     Letter     `gorm:"foreignKey:LetterID"`
	Manuscript Manuscript `gorm:"foreignKey:ManuscriptID"`
}
