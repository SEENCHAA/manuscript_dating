package main

import (
	"context" // Необходим для MinIO Context
	"log"     // Для логгирования MinIO-ошибок

	"lab31/internal/app/config"
	"lab31/internal/app/dsn"
	"lab31/internal/app/handler"
	"lab31/internal/app/redis"
	"lab31/internal/app/repository"
	"lab31/internal/pkg"

	_ "lab31/docs"

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

// @title Manuscript API
// @version 1.0
// @description API для работы с рукописями, письмами и пользователями.
// @contact.name API Support
// @contact.email support@example.com
// @license.name AS IS (NO WARRANTY)
// @host localhost:8081
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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

	// --- Инициализация Context, используем его для MinIO и Redis ---
	ctx := context.Background()

	// --- Инициализация MinIO (переменная ctx уже объявлена, используем только =) ---
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
	// --- Конец блока MinIO ---

	// --- Инициализация Redis (переменная ctx уже объявлена, используем только =) ---
	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		logrus.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()
	// --- Конец блока Redis ---

	// 3. Setup Router and Controller
	router := gin.Default()
	// ПЕРЕДАЧА MinIO-КЛИЕНТА, КОНФИГА И REDIS-КЛИЕНТА В КОНТРОЛЛЕР
	controller := handler.NewApplicationController(model, minioClient, cfg, redisClient)

	// 4. Create and Run Application
	app := pkg.NewApp(cfg, router, controller)
	app.RunApp()
}
