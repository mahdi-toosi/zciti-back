package user

import (
	"github.com/bangadam/go-fiber-starter/app/module/user/controller"
	"github.com/bangadam/go-fiber-starter/app/module/user/repository"
	"github.com/bangadam/go-fiber-starter/app/module/user/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Router struct {
	App        fiber.Router
	Controller *controller.Controller
}

func (_i *Router) RegisterRoutes() {
	// define controllers
	c := _i.Controller.RestController

	// define routes
	_i.App.Route("/api/v1", func(router fiber.Router) {
		router.Get("/", c.Index)
		router.Get("/:id", c.Show)
		router.Post("/", c.Store)
		router.Put("/:id", c.Update)
		router.Delete("/:id", c.Delete)
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
