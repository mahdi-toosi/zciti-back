package controller

import (
	"go-fiber-starter/app/module/comment/service"
	"go-fiber-starter/app/module/post/repository"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, r repository.IRepository) *Controller {
	return &Controller{
		RestController(s, r),
	}
}
