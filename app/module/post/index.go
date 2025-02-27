package post

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/post/controller"
	"go-fiber-starter/app/module/post/repository"
	"go-fiber-starter/app/module/post/service"
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
	_i.App.Route("/v1/business/:businessID/posts", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PReadAll), c.Index)
		router.Get("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PReadSingle), c.Show)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PCreate), c.Store)
		router.Post("/:id/delete-taxonomies", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PCreate), c.DeleteTaxonomies)
		router.Post("/:id/insert-taxonomies", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PCreate), c.InsertTaxonomies)
		router.Put("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DPost, mdl.PDelete), c.Delete)
	})

	_i.App.Route("/v1/user/business/:businessID/posts", func(router fiber.Router) {
		router.Get("/", mdl.ForUser, c.Index)
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
