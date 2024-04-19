package main

import (
	"github.com/efectn/fx-zerolog"
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
		// database
		fx.Provide(database.NewDatabase),

		// start application
		fx.Invoke(bootstrap.Seeder),

		// define logger
		fx.WithLogger(fxzerolog.Init()),
	).Run()
}
