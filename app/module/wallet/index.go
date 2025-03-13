package wallet

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/wallet/controller"
	"go-fiber-starter/app/module/wallet/repository"
	"go-fiber-starter/app/module/wallet/service"
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
	_i.App.Route("/v1/wallets", func(router fiber.Router) {
		router.Get("/wallet", mdl.Protected(cfg), c.Show)
		// router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DNotification, mdl.PReadAll), c.Index)
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
