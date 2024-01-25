package controller

import (
	businessService "go-fiber-starter/app/module/business/service"
	"go-fiber-starter/app/module/message/service"
	messageRoomService "go-fiber-starter/app/module/messageRoom/service"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, ms messageRoomService.IService, bs businessService.IService) *Controller {
	return &Controller{
		RestController(s, ms, bs),
	}
}
