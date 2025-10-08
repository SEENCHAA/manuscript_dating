package handler

import (
	"lab12/internal/app/repository"
)

type ApplicationController struct {
	ApplicationModel *repository.ApplicationModel
}

func NewApplicationController(r *repository.ApplicationModel) *ApplicationController {
	return &ApplicationController{
		ApplicationModel: r,
	}
}
