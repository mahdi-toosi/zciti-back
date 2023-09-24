package user

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/user/controller"
	"go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/app/module/user/service"
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
	_i.App.Route("/api/v1/users", func(router fiber.Router) {
		router.Get("/", middleware.Protected(), middleware.Permission(middleware.DUser, middleware.PReadAll), c.Index)
		router.Get("/:id", middleware.Protected(), middleware.Permission(middleware.DUser, middleware.PReadSingle), c.Show)
		router.Post("/", middleware.Protected(), middleware.Permission(middleware.DUser, middleware.PCreate), c.Store)
		router.Put("/:id", middleware.Protected(), middleware.Permission(middleware.DUser, middleware.PUpdate), c.Update)
		router.Delete("/:id", middleware.Protected(), middleware.Permission(middleware.DUser, middleware.PDelete), c.Delete)
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
