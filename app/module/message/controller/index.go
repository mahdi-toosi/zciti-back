package controller

import (
	"go-fiber-starter/app/module/message/service"
	messageRoomService "go-fiber-starter/app/module/messageRoom/service"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, ms messageRoomService.IService) *Controller {
	return &Controller{
		RestController(s, ms),
	}
}
