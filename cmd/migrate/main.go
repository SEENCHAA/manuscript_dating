package main

import (
	"lab1/internal/app/ds"  // <-- ИСПРАВЛЕНО
	"lab1/internal/app/dsn" // <-- ИСПРАВЛЕНО

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Letter{},
		&ds.Manuscript{},
		&ds.ManuscriptLetter{},
	)
	if err != nil {
		panic("can't migrate db: " + err.Error())
	}
}
