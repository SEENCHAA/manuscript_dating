package repository

import (
	"lab12/internal/app/ds"

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
	// Миграции
	if err := db.AutoMigrate(
		&ds.User{},
		&ds.Letter{},
		&ds.Manuscript{},
		&ds.ManuscriptLetter{},
	); err != nil {
		return nil, err
	}
	return &ApplicationModel{
		db: db,
	}, nil
}
