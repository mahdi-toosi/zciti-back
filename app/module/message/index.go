package message

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/message/controller"
	"go-fiber-starter/app/module/message/repository"
	"go-fiber-starter/app/module/message/service"
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
	_i.App.Route("/api/v1/messages", func(router fiber.Router) {
		router.Get("/:businessID", mdl.Protected(cfg), mdl.Permission(mdl.DMessage, mdl.PReadAll), c.Index)
		router.Post("/", mdl.Protected(cfg), mdl.Permission(mdl.DMessage, mdl.PCreate), c.Store)
		router.Put("/:id", mdl.Protected(cfg), mdl.Permission(mdl.DMessage, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(cfg), mdl.Permission(mdl.DMessage, mdl.PDelete), c.Delete)
	})

	//_i.App.Route("/api/v1/message-rooms", func(router fiber.Router) {
	//	router.Get("/", mdl.Protected(cfg), mdl.Permission(mdl.DMessageRoom, mdl.PReadAll), c.IndexMessageRooms)
	//	router.Post("/", mdl.Protected(cfg), mdl.Permission(mdl.DMessageRoom, mdl.PCreate), c.StoreMessageRoom)
	//	router.Delete("/:id", mdl.Protected(cfg), mdl.Permission(mdl.DMessageRoom, mdl.PDelete), c.DeleteMessageRoom)
	//})
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
