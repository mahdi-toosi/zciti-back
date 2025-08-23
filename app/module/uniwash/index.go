package uniwash

import (
	"github.com/gofiber/fiber/v2"
	mdl "go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/uniwash/controller"
	"go-fiber-starter/app/module/uniwash/repository"
	"go-fiber-starter/app/module/uniwash/service"
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
	_i.App.Route("/v1/business/:businessID/uni-wash", func(router fiber.Router) {
		router.Post("/send-command", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PReadAll), c.SendCommand)
		router.Post("/send-device-is-off-msg-to-user/:reservationID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PReadAll), c.SendDeviceIsOffMsgToUser)
		router.Get("/check-last-command-status/:reservationID", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PReadAll), c.CheckLastCommandStatus)
		router.Get("/device/reservation-options", mdl.Protected(cfg), mdl.BusinessPermission(mdl.DProduct, mdl.PReadAll), c.GetReservationOptions)
	})

	_i.App.Route("/v1/user/business/:businessID/uni-wash", func(router fiber.Router) {
		router.Post("/send-command", mdl.Protected(cfg), mdl.ForUser, c.SendCommand)
		router.Get("/reserved-machines", mdl.Protected(cfg), mdl.ForUser, c.IndexReservedMachines)
	})
}

func newRouter(fiber *fiber.App, controller *controller.Controller) *Router {
	return &Router{
		App:        fiber,
		Controller: controller,
	}
}

var Module = fx.Options(
	fx.Provide(service.Service),

	fx.Provide(repository.Repository),

	fx.Provide(controller.Controllers),

	fx.Provide(newRouter),
)
