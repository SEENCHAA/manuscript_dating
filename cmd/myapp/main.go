package main

import (
	"lab12/internal/app/config"
	"lab12/internal/app/dsn"
	"lab12/internal/app/handler"
	"lab12/internal/app/repository"
	"lab12/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()

	rep, err := repository.NewApplicationModel(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	hand := handler.NewApplicationController(rep)

	router := gin.Default()

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
