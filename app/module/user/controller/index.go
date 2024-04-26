package controller

import (
	bService "go-fiber-starter/app/module/business/service"
	"go-fiber-starter/app/module/user/service"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, b bService.IService) *Controller {
	return &Controller{
		RestController(s, b),
	}
}
