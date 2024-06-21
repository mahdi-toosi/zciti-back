package business

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/business/controller"
	"go-fiber-starter/app/module/business/repository"
	"go-fiber-starter/app/module/business/service"
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
	_i.App.Route("/v1/businesses", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PReadAll), c.Index)
		router.Get("/types", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PReadSingle), c.Types)
		router.Get("/:businessID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PReadSingle), c.Show)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PCreate), c.Store)
		router.Put("/:businessID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PUpdate), c.Update)
		router.Delete("/:businessID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PDelete), c.Delete)
	})

	_i.App.Route("/v1/user/businesses", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), c.MyBusinesses)
		router.Get("/:businessID", c.Show)
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
