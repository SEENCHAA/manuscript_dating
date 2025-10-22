package pkg

import (
	"fmt"

	"lab1/internal/app/config"
	"lab1/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.ApplicationController
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.ApplicationController) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	a.Handler.RegisterAPI(a.Router)

	for _, ri := range a.Router.Routes() {
		logrus.Printf("Route: %s %s\n", ri.Method, ri.Path)
	}

	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatalf("ошибка запуска сервера: %v", err)
	}

	logrus.Info("Server down")
}
