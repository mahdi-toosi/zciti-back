package bootstrap

import (
	"context"
	"go-fiber-starter/internal"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/router"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/response"
	"go.uber.org/fx"
)

func NewFiber(cfg *config.Config, baleBot *internal.BaleBot) *fiber.App {
	var errHandler = func(c *fiber.Ctx, _error error) error {
		if err := response.ErrorHandler(c, _error, baleBot); err != nil {
			return err
		}
		return nil
	}

	// setup
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          errHandler,
		ServerHeader:          cfg.App.Name,
		AppName:               cfg.App.Name,
		Prefork:               cfg.App.Prefork,
		EnablePrintRoutes:     cfg.App.PrintRoutes,
		IdleTimeout:           cfg.App.IdleTimeout * time.Second,
	})

	// pass production config to check it
	response.IsProduction = cfg.App.Production

	return app
}

func Start(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	fiber *fiber.App,
	router *router.Router,
	middlewares *middleware.Middleware,
	database *database.Database,
	log zerolog.Logger,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {

				// Register middlewares & routes
				middlewares.Register()
				router.Register()

				// Custom Startup Messages

				log.Info().Msg(fiber.Config().AppName + " is running at the moment!")

				// Debug information
				if !cfg.App.Production {
					//prefork := "Enabled"
					//procs := runtime.GOMAXPROCS(0)
					//if !cfg.App.Prefork {
					//	procs = 1
					//	prefork = "Disabled"
					//}

					//log.Info().Msgf("Version: %s", "-")
					//log.Info().Msgf("Serve On: %s", cfg.ParseAddress())
					//log.Info().Msgf("Prefork: %s", prefork)
					//log.Info().Msgf("Handlers: %d", fiber.HandlersCount())
					//log.Info().Msgf("Processes: %d", procs)
					log.Info().Msgf("PID: %d", os.Getpid())
				}

				// Listen the app (with TLS Support)
				//if cfg.App.TLS.Enable {
				//	log.Info().Msg("TLS support was enabled.")
				//
				//	if err := fiber.ListenTLS(cfg.App.Port, cfg.App.TLS.CertFile, cfg.App.TLS.KeyFile); err != nil {
				//		log.Error().Err(err).Msg("An unknown error occurred when to run server!")
				//	}
				//}

				go func() {
					if err := fiber.Listen(":" + cfg.App.Port); err != nil {
						log.Error().Err(err).Msg("An unknown error occurred when to run server!")
					}
				}()

				//redis.ConnectRedis()
				database.ConnectDatabase()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().Msg("Shutting down the app...")
				if err := fiber.Shutdown(); err != nil {
					log.Panic().Err(err).Msg("")
				}

				log.Info().Msg("Running cleanup tasks...")
				log.Info().Msg("1- Shutdown the database")
				database.ShutdownDatabase()
				log.Info().Msgf("%s was successful shutdown.", cfg.App.Name)
				log.Info().Msg("\u001b[96m see you again👋 \u001b[0m")

				return nil
			},
		},
	)
}
