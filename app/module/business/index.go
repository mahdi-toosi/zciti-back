package business

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/business/controller"
	"go-fiber-starter/app/module/business/repository"
	"go-fiber-starter/app/module/business/service"
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
	_i.App.Route("/api/v1/businesses", func(router fiber.Router) {
		router.Get("/", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PReadAll), c.Index)
		router.Get("/types", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PReadSingle), c.Types)
		router.Get("/:id", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PReadSingle), c.Show)
		router.Get("/:id/users", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PReadSingle), c.Users)
		router.Post("/:businessID/users/:userID", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PReadSingle), c.InsertUser)
		router.Delete("/:businessID/users/:userID", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PDelete), c.DeleteUser)
		router.Post("/", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PCreate), c.Store)
		router.Put("/:id", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(), mdl.Permission(mdl.DBusiness, mdl.PDelete), c.Delete)
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
