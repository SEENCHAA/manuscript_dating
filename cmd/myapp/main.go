package main

import (
	"context" // Необходим для MinIO Context
	"log"     // Для логгирования MinIO-ошибок

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"lab1/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	// Импорты для MinIO
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Константы MinIO (лучше вынести в config, но здесь для примера)
const (
	MinioEndpoint  = "127.0.0.1:9000"
	MinioAccessKey = "minioadmin"
	MinioSecretKey = "minioadmin"
	MinioBucket    = "manuscripts"
)

func main() {
	// 1. Load config
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	// 2. Prepare DSN and Model (DB Connection)
	dsnString := dsn.FromEnv()
	if dsnString == "" {
		logrus.Fatal("DB connection string is empty. Check your .env file.")
	}
	model, err := repository.NewApplicationModel(dsnString)
	if err != nil {
		logrus.Fatalf("failed to connect database: %v", err)
	}

	// --- НОВЫЙ БЛОК: Инициализация MinIO ---

	ctx := context.Background()

	minioClient, err := minio.New(MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
		Secure: false, // Используйте true, если MinIO настроен на HTTPS
	})
	if err != nil {
		log.Fatalf("Ошибка инициализации MinIO-клиента: %v", err)
	}

	// Проверка и создание бакета, если он не существует
	found, err := minioClient.BucketExists(ctx, MinioBucket)
	if err != nil {
		log.Fatalf("Ошибка проверки существования бакета: %v", err)
	}
	if !found {
		log.Printf("Создание MinIO бакета '%s'", MinioBucket)
		err = minioClient.MakeBucket(ctx, MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Не удалось создать MinIO бакет: %v", err)
		}
	}
	// --- КОНЕЦ НОВОГО БЛОКА ---

	// 3. Setup Router and Controller
	router := gin.Default()
	// ПЕРЕДАЧА MinIO-КЛИЕНТА В КОНТРОЛЛЕР
	appController := handler.NewApplicationController(model, minioClient)

	// 4. Create and Run Application
	app := pkg.NewApp(cfg, router, appController)
	app.RunApp()
}
