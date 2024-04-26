package product

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	postController "go-fiber-starter/app/module/post/controller"
	"go-fiber-starter/app/module/product/controller"
	"go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/product/service"
	"go-fiber-starter/utils/config"
	"go.uber.org/fx"
)

type Router struct {
	App            fiber.Router
	Controller     *controller.Controller
	PostController *postController.Controller
}

func (_i *Router) RegisterRoutes(cfg *config.Config) {
	// define controllers
	c := _i.Controller.RestController
	pc := _i.PostController.RestController

	// define routes
	_i.App.Route("/v1/business/:businessID/products", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PReadAll), c.Index)
		router.Get("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PReadSingle), c.Show)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PCreate), c.Store)
		router.Post("/:id/delete-taxonomies", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PCreate), pc.DeleteTaxonomies)
		router.Post("/:id/insert-taxonomies", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PCreate), pc.InsertTaxonomies)
		router.Put("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PDelete), c.Delete)
	})
}

func newRouter(fiber *fiber.App, controller *controller.Controller, PostController *postController.Controller) *Router {
	return &Router{
		App:            fiber,
		Controller:     controller,
		PostController: PostController,
	}
}

var Module = fx.Options(
	fx.Provide(repository.Repository),

	fx.Provide(service.Service),

	fx.Provide(controller.Controllers),

	fx.Provide(newRouter),
)
