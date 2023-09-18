package controller

import "github.com/bangadam/go-fiber-starter/app/module/auth/service"

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService) *Controller {
	return &Controller{
		RestController(s),
	}
}
