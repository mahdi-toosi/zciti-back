package transaction

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/transaction/controller"
	"go-fiber-starter/app/module/transaction/repository"
	"go-fiber-starter/app/module/transaction/service"
	"go-fiber-starter/utils/config"
	"go.uber.org/fx"
)

type Router struct {
	App        fiber.Router
	Controller *controller.Controller
}

func (_i *Router) RegisterRoutes(cfg *config.Config) {
	// define controllers
	c := _i.Controller.RestController

	// define routes
	_i.App.Route("/v1/wallets/:walletID/transactions", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), c.Index)
	})
}

func newRouter(fiber *fiber.App, controller *controller.Controller) *Router {
	return &Router{
		App:        fiber,
		Controller: controller,
	}
}

var Module = fx.Options(
	fx.Provide(repository.Repository),

	fx.Provide(service.Service),

	fx.Provide(controller.Controllers),

	fx.Provide(newRouter),
)
