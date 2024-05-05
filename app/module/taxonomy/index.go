package taxonomy

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/taxonomy/controller"
	"go-fiber-starter/app/module/taxonomy/repository"
	"go-fiber-starter/app/module/taxonomy/service"
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
	_i.App.Route("/v1/business/:businessID/taxonomies", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DTaxonomy, mdl.PReadAll), c.Index)
		router.Get("/search", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DTaxonomy, mdl.PReadAll), c.Search)
		router.Get("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DTaxonomy, mdl.PReadSingle), c.Show)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DTaxonomy, mdl.PCreate), c.Store)
		router.Put("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DTaxonomy, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DTaxonomy, mdl.PDelete), c.Delete)
	})

	_i.App.Route("/v1/user/business/:businessID/taxonomies", func(router fiber.Router) {
		router.Get("/", mdl.ForUser, c.Search)
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
