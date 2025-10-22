package handler

import (
	"lab1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7" // Необходимый импорт
)

type ApplicationController struct {
	ApplicationModel *repository.ApplicationModel
	MinioClient      *minio.Client // <-- Поле для MinIO-клиента
}

// NewApplicationController теперь принимает MinIO-клиент
func NewApplicationController(model *repository.ApplicationModel, minioClient *minio.Client) *ApplicationController {
	return &ApplicationController{
		ApplicationModel: model,
		MinioClient:      minioClient, // <-- Корректное присвоение
	}
}

// RegisterRoutes вызывает RegisterAPI для регистрации всех маршрутов.
func (h *ApplicationController) RegisterRoutes(router *gin.Engine) {
	h.RegisterAPI(router)
}

// RegisterAPI регистрирует все API маршруты с использованием методов контроллера.
func (h *ApplicationController) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		// 1-5. USERS (5 маршрутов)
		users := api.Group("/users")
		{
			users.POST("/register", h.ApiRegisterUser) // 1. POST регистрация
			users.POST("/login", h.ApiLoginUser)       // 2. POST аутентификация
			users.POST("/logout", h.ApiLogoutUser)     // 3. POST деавторизация
			users.GET("/me", h.ApiGetUser)             // 4. GET мои данные
			users.PUT("/me", h.ApiUpdateUser)          // 5. PUT изменение
		}

		// 6-12. LETTERS (7 маршрутов)
		letters := api.Group("/letters")
		{
			letters.GET("", h.ApiGetLetters)                       // 6. GET список
			letters.GET("/:id", h.ApiGetLetter)                    // 7. GET одна запись
			letters.POST("", h.ApiCreateLetter)                    // 8. POST создание
			letters.PUT("/:id", h.ApiUpdateLetter)                 // 9. PUT изменение
			letters.DELETE("/:id", h.ApiDeleteLetter)              // 10. DELETE удаление
			letters.POST("/:id/image", h.ApiUploadImage)           // 11. POST добавление изображения (MinIO)
			letters.POST("/:id/manuscript", h.ApiAddLetterToDraft) // 12. POST добавление в заявку-черновик
		}

		// 13-21. MANUSCRIPTS & M-M (9 маршрутов)
		manuscripts := api.Group("/manuscripts")
		{
			manuscripts.GET("/basket", h.ApiGetBasketStatus) // 13. GET иконки корзины (заменено с /draft)
			manuscripts.GET("", h.ApiGetManuscripts)         // 14. GET список
			manuscripts.GET("/:id", h.ApiGetManuscript)      // 15. GET одна запись

			manuscripts.PUT("/:id", h.ApiUpdateManuscript)              // 16. PUT изменения полей заявки по теме
			manuscripts.PUT("/:id/submit", h.ApiSubmitManuscript)       // 17. PUT сформировать создателем
			manuscripts.PUT("/:id/moderation", h.ApiModerateManuscript) // 18. PUT завершить/отклонить модератором (ОБЪЕДИНЕН)

			manuscripts.DELETE("/:id", h.ApiDeleteManuscript) // 19. DELETE заявки (статус -> deleted)

			manuscripts.PUT("/:id/letters/:lid", h.ApiUpdateMMQuantity) // 20. PUT изменение количества/порядка/значения в м-м
			manuscripts.DELETE("/:id/letters/:lid", h.ApiDeleteMM)      // 21. DELETE удаление услуги из заявки (DELETE м-м)
		}
	}
}
