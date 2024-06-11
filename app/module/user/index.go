package user

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/user/controller"
	"go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/app/module/user/service"
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
	_i.App.Route("/v1/users", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DUser, mdl.PReadAll), c.Index)
		router.Get("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DUser, mdl.PReadSingle), c.Show)
		router.Post("/", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DUser, mdl.PCreate), c.Store)
		router.Put("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DUser, mdl.PUpdate), c.Update)
		router.Delete("/:id", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DUser, mdl.PDelete), c.Delete)
	})

	_i.App.Route("/v1/business/:businessID", func(router fiber.Router) {
		router.Get("/users", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PReadSingle), c.BusinessUsers)
		router.Post("/users/:userID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PReadSingle), c.InsertUser)
		router.Post("/users-add-role", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PUpdate), c.BusinessUsersAddRole)
		router.Delete("/users/:userID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DBusiness, mdl.PDelete), c.DeleteUser)
	})

	// TODO move it to orders module
	_i.App.Route("/v1/user/orders", func(router fiber.Router) {
		router.Get("/", mdl.Protected(cfg), c.Orders)
		router.Post("/", mdl.Protected(cfg), c.OrderStore)
		router.Get("/status", c.OrderStatus)
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
