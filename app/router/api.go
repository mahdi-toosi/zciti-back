package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go-fiber-starter/app/module/asset"
	"go-fiber-starter/app/module/auth"
	"go-fiber-starter/app/module/business"
	"go-fiber-starter/app/module/comment"
	"go-fiber-starter/app/module/message"
	"go-fiber-starter/app/module/messageRoom"
	"go-fiber-starter/app/module/notification"
	"go-fiber-starter/app/module/notificationTemplate"
	"go-fiber-starter/app/module/post"
	"go-fiber-starter/app/module/product"
	"go-fiber-starter/app/module/taxonomy"
	"go-fiber-starter/app/module/user"
	"go-fiber-starter/utils/config"
)

type Router struct {
	App fiber.Router
	Cfg *config.Config

	AuthRouter                 *auth.Router
	UserRouter                 *user.Router
	PostRouter                 *post.Router
	AssetRouter                *asset.Router
	ProductRouter              *product.Router
	MessageRouter              *message.Router
	CommentRouter              *comment.Router
	TaxonomyRouter             *taxonomy.Router
	BusinessRouter             *business.Router
	MessageRoomRouter          *messageRoom.Router
	NotificationRouter         *notification.Router
	NotificationTemplateRouter *notificationtemplate.Router
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,

	authRouter *auth.Router,
	userRouter *user.Router,
	postRouter *post.Router,
	assetRouter *asset.Router,
	productRouter *product.Router,
	commentRouter *comment.Router,
	messageRouter *message.Router,
	taxonomyRouter *taxonomy.Router,
	businessRouter *business.Router,
	messageRoomRouter *messageRoom.Router,
	notificationRouter *notification.Router,
	notificationTemplateRouter *notificationtemplate.Router,
) *Router {
	return &Router{
		App: fiber,
		Cfg: cfg,

		AuthRouter:                 authRouter,
		UserRouter:                 userRouter,
		PostRouter:                 postRouter,
		AssetRouter:                assetRouter,
		ProductRouter:              productRouter,
		CommentRouter:              commentRouter,
		MessageRouter:              messageRouter,
		TaxonomyRouter:             taxonomyRouter,
		BusinessRouter:             businessRouter,
		MessageRoomRouter:          messageRoomRouter,
		NotificationRouter:         notificationRouter,
		NotificationTemplateRouter: notificationTemplateRouter,
	}
}

// Register routes
func (r *Router) Register() { // Register routes of modules
	r.AuthRouter.RegisterRoutes()
	r.UserRouter.RegisterRoutes(r.Cfg)
	r.PostRouter.RegisterRoutes(r.Cfg)
	r.AssetRouter.RegisterRoutes(r.Cfg)
	r.ProductRouter.RegisterRoutes(r.Cfg)
	r.CommentRouter.RegisterRoutes(r.Cfg)
	r.MessageRouter.RegisterRoutes(r.Cfg)
	r.TaxonomyRouter.RegisterRoutes(r.Cfg)
	r.BusinessRouter.RegisterRoutes(r.Cfg)
	r.MessageRoomRouter.RegisterRoutes(r.Cfg)
	r.NotificationRouter.RegisterRoutes(r.Cfg)
	r.NotificationTemplateRouter.RegisterRoutes(r.Cfg)

	// Swagger Documentation
	r.App.Get("/swagger/*", swagger.HandlerDefault)
}
