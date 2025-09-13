package api

import (
	"lab12/internal/app/handler"
	"lab12/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Start")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	// Загружаем шаблоны
	r.LoadHTMLGlob("templates/*.html")

	// Раздаём статику
	r.Static("/resources", "./resources")

	// Маршруты
	r.GET("/", h.GetFeatures)
	r.POST("/", h.GetFeatures)

	r.GET("/feature/:id", h.GetFeature)  // правильный маршрут
	r.POST("/feature/:id", h.GetFeature) // POST тоже работает

	r.GET("/order", h.Order)
	r.GET("/order/delete", h.Delete)

	// Выводим все зарегистрированные маршруты для проверки
	for _, ri := range r.Routes() {
		log.Printf("Route: %s %s\n", ri.Method, ri.Path)
	}

	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("ошибка запуска сервера: %v", err)
	}

	log.Println("Server down")
}
