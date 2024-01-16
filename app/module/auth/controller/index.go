package controller

import (
	"go-fiber-starter/app/module/auth/service"
	"go-fiber-starter/utils/config"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, c *config.Config) *Controller {
	return &Controller{
		RestController(s, c),
	}
}
