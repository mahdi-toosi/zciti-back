package controller

import "go-fiber-starter/app/module/order/service"

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService) *Controller {
	return &Controller{
		RestController(s),
	}
}
