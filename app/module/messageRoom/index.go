package messageRoom

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/messageRoom/controller"
	"go-fiber-starter/app/module/messageRoom/repository"
	"go-fiber-starter/app/module/messageRoom/service"
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
	_i.App.Route("/api/v1/message-rooms", func(router fiber.Router) {
		router.Get("/", mdl.Protected(), mdl.Permission(mdl.DMessageRoom, mdl.PReadAll), c.Index)
		router.Delete("/:id", mdl.Protected(), mdl.Permission(mdl.DMessageRoom, mdl.PDelete), c.Delete)
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
