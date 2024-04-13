package comment

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/comment/controller"
	"go-fiber-starter/app/module/comment/repository"
	"go-fiber-starter/app/module/comment/service"
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
	_i.App.Route("/v1/business/:businessID/posts/:postID/comments", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DComment, mdl.PReadAll), c.Index)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DComment, mdl.PCreate), c.Store)
		router.Put("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DComment, mdl.PUpdate), c.Update)
		router.Put("/:id/status", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DComment, mdl.PUpdate), c.UpdateStatus)
		//router.Delete("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DComment, mdl.PDelete), c.Delete)
		//router.Get("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DComment, mdl.PReadSingle), c.Show)
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
