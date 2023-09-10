package router

import (
	"github.com/bangadam/go-fiber-starter/app/module/auth"
	"github.com/bangadam/go-fiber-starter/app/module/user"
	"github.com/bangadam/go-fiber-starter/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Router struct {
	App        fiber.Router
	Cfg        *config.Config
	AuthRouter *auth.Router
	UserRouter *user.Router
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,
	userRouter *user.Router,
	authRouter *auth.Router) *Router {
	return &Router{
		App:        fiber,
		Cfg:        cfg,
		AuthRouter: authRouter,
		UserRouter: userRouter,
	}
}

// Register routes
func (r *Router) Register() {

	// Register routes of modules
	r.UserRouter.RegisterRoutes()
	r.AuthRouter.RegisterRoutes()

	// Swagger Documentation
	r.App.Get("/swagger/*", swagger.HandlerDefault)
}
