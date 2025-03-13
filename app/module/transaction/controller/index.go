package controller

import (
	"go-fiber-starter/app/module/transaction/service"
	walletService "go-fiber-starter/app/module/wallet/service"
)

type Controller struct {
	RestController IRestController
}

func Controllers(s service.IService, w walletService.IService) *Controller {
	return &Controller{
		RestController(s, w),
	}
}
