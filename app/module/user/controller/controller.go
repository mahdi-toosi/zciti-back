package controller

import "github.com/bangadam/go-fiber-starter/app/module/user/service"

type Controller struct {
	User IRestController
}

func Controllers(s service.IService) *Controller {
	return &Controller{
		User: RestController(s),
	}
}
