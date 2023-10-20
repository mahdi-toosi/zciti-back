package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go-fiber-starter/app/module/auth"
	"go-fiber-starter/app/module/notification"
	"go-fiber-starter/app/module/notificationTemplate"
	"go-fiber-starter/app/module/post"
	"go-fiber-starter/app/module/user"
	"go-fiber-starter/utils/config"
)

type Router struct {
	App fiber.Router
	Cfg *config.Config

	AuthRouter                 *auth.Router
	UserRouter                 *user.Router
	PostRouter                 *post.Router
	NotificationRouter         *notification.Router
	NotificationTemplateRouter *notificationtemplate.Router
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,

	authRouter *auth.Router,
	userRouter *user.Router,
	postRouter *post.Router,
	notificationRouter *notification.Router,
	NotificationTemplateRouter *notificationtemplate.Router,
) *Router {
	return &Router{
		App: fiber,
		Cfg: cfg,

		AuthRouter:                 authRouter,
		UserRouter:                 userRouter,
		PostRouter:                 postRouter,
		NotificationRouter:         notificationRouter,
		NotificationTemplateRouter: NotificationTemplateRouter,
	}
}

// Register routes
func (r *Router) Register() { // Register routes of modules
	r.UserRouter.RegisterRoutes()
	r.AuthRouter.RegisterRoutes()
	r.PostRouter.RegisterRoutes()
	r.NotificationRouter.RegisterRoutes()
	r.NotificationTemplateRouter.RegisterRoutes()

	// Swagger Documentation
	r.App.Get("/swagger/*", swagger.HandlerDefault)
}
