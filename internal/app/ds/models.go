package ds

import (
	"lab31/internal/app/role"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Role     role.Role `gorm:"type:integer;not null;default:0"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   uint      `json:"user_id"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Role     role.Role `json:"role"`
}

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

type Manuscript struct {
	ID               uint      `gorm:"primaryKey"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	Status           string    `gorm:"type:varchar(20);not null;default:'draft'"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	SubmittedAt      *time.Time
	FinishedAt       *time.Time
	ModeratorID      *uint
	CalculatedPeriod string `gorm:"type:varchar(50)"`
	ManuscriptText   string `gorm:"type:text" json:"manuscript_text"`

	User      User               `gorm:"foreignKey:UserID"`
	Moderator User               `gorm:"foreignKey:ModeratorID"`
	Letters   []ManuscriptLetter `gorm:"foreignKey:ManuscriptID"`
}

type ManuscriptLetter struct {
	ManuscriptID uint `gorm:"primaryKey"`
	LetterID     uint `gorm:"primaryKey"`
	Quantity     int  `gorm:"not null;default:1"`

	Manuscript Manuscript `gorm:"foreignKey:ManuscriptID"`
	Letter     Letter     `gorm:"foreignKey:LetterID"`
}
