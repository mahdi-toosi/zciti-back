package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go-fiber-starter/app/module/asset"
	"go-fiber-starter/app/module/auth"
	"go-fiber-starter/app/module/business"
	"go-fiber-starter/app/module/comment"
	"go-fiber-starter/app/module/coupon"
	"go-fiber-starter/app/module/notification"
	"go-fiber-starter/app/module/notificationTemplate"
	"go-fiber-starter/app/module/order"
	"go-fiber-starter/app/module/post"
	"go-fiber-starter/app/module/product"
	"go-fiber-starter/app/module/reservation"
	"go-fiber-starter/app/module/taxonomy"
	"go-fiber-starter/app/module/transaction"
	"go-fiber-starter/app/module/uniwash"
	"go-fiber-starter/app/module/user"
	"go-fiber-starter/app/module/wallet"
	"go-fiber-starter/utils/config"
)

type Router struct {
	App fiber.Router
	Cfg *config.Config

	AuthRouter                 *auth.Router
	UserRouter                 *user.Router
	PostRouter                 *post.Router
	OrderRouter                *order.Router
	AssetRouter                *asset.Router
	WalletRouter               *wallet.Router
	CouponRouter               *coupon.Router
	UniWashRouter              *uniwash.Router
	ProductRouter              *product.Router
	CommentRouter              *comment.Router
	TaxonomyRouter             *taxonomy.Router
	BusinessRouter             *business.Router
	TransactionRouter          *transaction.Router
	ReservationsRouter         *reservation.Router
	NotificationRouter         *notification.Router
	NotificationTemplateRouter *notificationtemplate.Router
	//MessageRouter              *message.Router
	//MessageRoomRouter          *messageRoom.Router
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,

	authRouter *auth.Router,
	userRouter *user.Router,
	postRouter *post.Router,
	orderRouter *order.Router,
	assetRouter *asset.Router,
	walletRouter *wallet.Router,
	couponRouter *coupon.Router,
	productRouter *product.Router,
	uniWashRouter *uniwash.Router,
	commentRouter *comment.Router,
	taxonomyRouter *taxonomy.Router,
	businessRouter *business.Router,
	transactionRouter *transaction.Router,
	reservationsRouter *reservation.Router,
	notificationRouter *notification.Router,
	notificationTemplateRouter *notificationtemplate.Router,
	// messageRouter *message.Router,
	// messageRoomRouter *messageRoom.Router,
) *Router {
	return &Router{
		App: fiber,
		Cfg: cfg,

		AuthRouter:    authRouter,
		UserRouter:    userRouter,
		PostRouter:    postRouter,
		OrderRouter:   orderRouter,
		AssetRouter:   assetRouter,
		WalletRouter:  walletRouter,
		CouponRouter:  couponRouter,
		ProductRouter: productRouter,
		UniWashRouter: uniWashRouter,
		CommentRouter: commentRouter,
		//MessageRouter:              messageRouter,
		TaxonomyRouter:     taxonomyRouter,
		BusinessRouter:     businessRouter,
		TransactionRouter:  transactionRouter,
		ReservationsRouter: reservationsRouter,
		//MessageRoomRouter:          messageRoomRouter,
		NotificationRouter:         notificationRouter,
		NotificationTemplateRouter: notificationTemplateRouter,
	}
}

// Register routes
func (r *Router) Register() { // Register routes of modules
	r.AuthRouter.RegisterRoutes()
	r.UserRouter.RegisterRoutes(r.Cfg)
	r.PostRouter.RegisterRoutes(r.Cfg)
	r.OrderRouter.RegisterRoutes(r.Cfg)
	r.AssetRouter.RegisterRoutes(r.Cfg)
	r.WalletRouter.RegisterRoutes(r.Cfg)
	r.CouponRouter.RegisterRoutes(r.Cfg)
	r.ProductRouter.RegisterRoutes(r.Cfg)
	r.UniWashRouter.RegisterRoutes(r.Cfg)
	r.CommentRouter.RegisterRoutes(r.Cfg)
	//r.MessageRouter.RegisterRoutes(r.Cfg)
	r.TaxonomyRouter.RegisterRoutes(r.Cfg)
	r.BusinessRouter.RegisterRoutes(r.Cfg)
	r.TransactionRouter.RegisterRoutes(r.Cfg)
	r.ReservationsRouter.RegisterRoutes(r.Cfg)
	//r.MessageRoomRouter.RegisterRoutes(r.Cfg)
	r.NotificationRouter.RegisterRoutes(r.Cfg)
	r.NotificationTemplateRouter.RegisterRoutes(r.Cfg)

	// Swagger Documentation
	r.App.Get("/swagger/*", swagger.HandlerDefault)
	r.App.Get("/health-check", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
