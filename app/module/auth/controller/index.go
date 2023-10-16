package controller

import (
	"go-fiber-starter/app/module/auth/service"
	usersRepo "go-fiber-starter/app/module/user/repository"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, u usersRepo.IRepository) *Controller {
	return &Controller{
		RestController(s, u),
	}
}
