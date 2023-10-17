package auth

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/auth/controller"
	"go-fiber-starter/app/module/auth/service"
	"go.uber.org/fx"
)

type Router struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewRouter(fiber *fiber.App, controller *controller.Controller) *Router {
	return &Router{
		App:        fiber,
		Controller: controller,
	}
}

func (_i *Router) RegisterRoutes() {
	c := _i.Controller.RestController

	_i.App.Route("/api/v1", func(router fiber.Router) {
		router.Post("/auth/login", c.Login)
		router.Post("/auth/register", c.Register)
		router.Post("/auth/send-otp", c.SendOtp)
		router.Post("/auth/reset-pass", c.ResetPass)
		//router.Post("/auth/forgot-pass", c.Login)
	})
}

var Module = fx.Options(
	fx.Provide(service.Service),

	fx.Provide(controller.Controllers),

	fx.Provide(NewRouter),
)
