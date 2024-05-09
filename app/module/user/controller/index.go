package controller

import (
	bService "go-fiber-starter/app/module/business/service"
	oService "go-fiber-starter/app/module/order/service"
	"go-fiber-starter/app/module/user/service"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, b bService.IService, o oService.IService) *Controller {
	return &Controller{
		RestController(s, b, o),
	}
}
