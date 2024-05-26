package controller

import (
	bService "go-fiber-starter/app/module/business/service"
	oService "go-fiber-starter/app/module/order/service"
	"go-fiber-starter/app/module/user/service"
	"go-fiber-starter/utils/config"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, b bService.IService, o oService.IService, config *config.Config) *Controller {
	return &Controller{
		RestController(s, b, o, config),
	}
}
