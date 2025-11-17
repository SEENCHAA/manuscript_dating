package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"lab31/internal/app/config"
	"lab31/internal/app/ds"
	"lab31/internal/app/redis"
	"lab31/internal/app/repository"
	"lab31/internal/app/role"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swaggo files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type ApplicationController struct {
	ApplicationModel *repository.ApplicationModel
	MinioClient      *minio.Client
	Config           *config.Config
	RedisClient      *redis.Client
}

func NewApplicationController(model *repository.ApplicationModel, minioClient *minio.Client, cfg *config.Config, redisClient *redis.Client) *ApplicationController {
	return &ApplicationController{
		ApplicationModel: model,
		MinioClient:      minioClient,
		Config:           cfg,
		RedisClient:      redisClient,
	}
}

func (h *ApplicationController) RegisterRoutes(router *gin.Engine) {
	h.RegisterAPI(router)
}

func (h *ApplicationController) RegisterAPI(router *gin.Engine) {
	adminOnly := h.WithAuthCheck(role.Admin)
	moderatorAndAdmin := h.WithAuthCheck(role.Admin, role.Manager)
	userAuth := h.WithAuthCheck(role.Admin, role.Manager, role.Buyer)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", h.ApiRegisterUser)
			users.POST("/login", h.ApiLoginUser)
			users.POST("/logout", userAuth, h.ApiLogoutUser)
			users.GET("/me", userAuth, h.ApiGetUser)
			users.PUT("/me", userAuth, h.ApiUpdateUser)
		}

		letters := api.Group("/letters")
		{
			letters.GET("", h.ApiGetLetters)
			letters.GET("/:id", h.ApiGetLetter)

			letters.POST("", moderatorAndAdmin, h.ApiCreateLetter)
			letters.PUT("/:id", moderatorAndAdmin, h.ApiUpdateLetter)
			letters.DELETE("/:id", adminOnly, h.ApiDeleteLetter)
			letters.POST("/:id/image", moderatorAndAdmin, h.ApiUploadImage)
			letters.POST("/:id/manuscript", userAuth, h.ApiAddLetterToDraft)
		}

		manuscripts := api.Group("/manuscripts")
		{
			manuscripts.GET("/basket", userAuth, h.ApiGetBasketStatus)
			manuscripts.GET("", userAuth, h.ApiGetManuscripts)
			manuscripts.GET("/:id", userAuth, h.ApiGetManuscript)

			manuscripts.PUT("/:id", userAuth, h.ApiUpdateManuscript)
			manuscripts.PUT("/:id/submit", moderatorAndAdmin, h.ApiSubmitManuscript)
			manuscripts.PUT("/:id/moderation", moderatorAndAdmin, h.ApiModerateManuscript)
			manuscripts.DELETE("/:id", userAuth, h.ApiDeleteManuscript)

			manuscripts.PUT("/:id/letters/:lid", userAuth, h.ApiUpdateMMQuantity)
			manuscripts.DELETE("/:id/letters/:lid", userAuth, h.ApiDeleteMM)
		}
	}
}

func (h *ApplicationController) WithAuthCheck(assignedRoles ...role.Role) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		authHeader := gCtx.GetHeader("Authorization")
		if authHeader == "" {
			gCtx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwtStr := strings.TrimPrefix(authHeader, "Bearer ")
		if jwtStr == authHeader {
			gCtx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 2. Проверка Blacklist
		_, err := h.RedisClient.GetBlacklistEntry(gCtx, jwtStr).Result()

		if err == nil {
			// ИСПРАВЛЕНИЕ: Если err == nil, ключ найден, токен в черном списке
			gCtx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !errors.Is(err, goredis.Nil) {
			// Внутренняя ошибка Redis (не "ключ не найден")
			logrus.Errorf("Redis internal error on blacklist check: %v", err)
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Если err == goredis.Nil, токен не в черном списке. Продолжаем.

		// 3. Парсинг и валидация токена
		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.JWT.Token), nil
		})
		if err != nil || !token.Valid { // Добавил проверку .Valid
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Println(err)
			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		gCtx.Set("user_id", myClaims.UserID)
		gCtx.Set("user_role", myClaims.Role)

		// 4. Проверка роли
		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Next()
				return
			}
		}

		// Роль не найдена
		gCtx.AbortWithStatus(http.StatusForbidden)
		logrus.Printf("role %s is not assigned in %s", myClaims.Role, assignedRoles)
		return
	}
}
