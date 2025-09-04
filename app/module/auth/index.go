package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go-fiber-starter/app/module/auth/controller"
	"go-fiber-starter/app/module/auth/service"
	"go.uber.org/fx"
	"time"
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

	_i.App.Route("/v1", func(router fiber.Router) {
		router.Post("/auth/login", c.Login)
		router.Post("/auth/register", c.Register)
		router.Post("/auth/reset-pass", c.ResetPass)
		router.Post("/auth/send-otp", limiter.New(limiter.Config{
			Max:        3,
			Expiration: 1 * time.Hour,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "تعداد درخواست شما به بیش از حد مجاز رسیده است.",
				})
			},
		}), c.SendOtp)
	})
}

var Module = fx.Options(
	fx.Provide(service.Service),

	fx.Provide(controller.Controllers),

	fx.Provide(NewRouter),
)
