package notificationtemplate

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/notificationTemplate/controller"
	"go-fiber-starter/app/module/notificationTemplate/repository"
	"go-fiber-starter/app/module/notificationTemplate/service"
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
	_i.App.Route("/api/v1/business/:businessID/notificationTemplates", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DNotificationTemplate, mdl.PReadAll), c.Index)
		router.Get("/keywords", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DNotificationTemplate, mdl.PReadSingle), c.Keywords)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DNotificationTemplate, mdl.PCreate), c.Store)
		router.Patch("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DNotificationTemplate, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DNotificationTemplate, mdl.PDelete), c.Delete)
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
