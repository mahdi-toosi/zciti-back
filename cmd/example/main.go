package main

import (
	"go-fiber-starter/app/module/auth"
	"go-fiber-starter/app/module/business"
	"go-fiber-starter/app/module/message"
	"go-fiber-starter/app/module/messageRoom"
	"go-fiber-starter/app/module/notification"
	"go-fiber-starter/app/module/notificationTemplate"
	"go-fiber-starter/app/module/post"
	"go-fiber-starter/app/module/user"

	"github.com/efectn/fx-zerolog"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/router"
	"go-fiber-starter/internal"
	"go-fiber-starter/internal/bootstrap"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/config"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
)

// @title                       Go Fiber Starter API Documentation
// @version                     1.0
// @description                 This is a sample API documentation.
// @termsOfService              http://swagger.io/terms/
// @contact.name                Mahdi Toosi
// @contact.email               mailmahditoosi@gmail.com
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @host                        localhost:8000/api/v1
// @schemes                     http
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 "Type 'Bearer {TOKEN}' to correctly set the API Key"
// @BasePath                    /
func main() {
	fx.New(
		/* provide patterns */
		// config
		fx.Provide(config.NewConfig),
		// logging
		fx.Provide(bootstrap.NewLogger),
		// fiber
		fx.Provide(bootstrap.NewFiber),
		// database
		fx.Provide(database.NewDatabase),
		// redis
		fx.Provide(bootstrap.NewRedis),
		// middleware
		fx.Provide(middleware.NewMiddleware),
		// router
		fx.Provide(router.NewRouter),
		// messageWay service
		fx.Provide(internal.NewMessageWay),

		// provide modules
		post.Module,
		user.Module,
		auth.Module,
		message.Module,
		business.Module,
		messageRoom.Module,
		notification.Module,
		notificationtemplate.Module,
		// End provide modules

		// start application
		fx.Invoke(bootstrap.Start),

		// define logger
		fx.WithLogger(fxzerolog.Init()),
	).Run()
}
