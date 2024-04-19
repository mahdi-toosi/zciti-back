package bootstrap

import (
	"context"
	"os"
	"runtime"
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

func NewFiber(cfg *config.Config) *fiber.App {
	// setup
	app := fiber.New(fiber.Config{
		ServerHeader:          cfg.App.Name,
		AppName:               cfg.App.Name,
		Prefork:               cfg.App.Prefork,
		ErrorHandler:          response.ErrorHandler,
		IdleTimeout:           cfg.App.IdleTimeout * time.Second,
		EnablePrintRoutes:     cfg.App.PrintRoutes,
		DisableStartupMessage: true,
	})

	// pass production config to check it
	response.IsProduction = cfg.App.Production

	return app
}

func Start(lifecycle fx.Lifecycle, cfg *config.Config, fiber *fiber.App, router *router.Router, middlewares *middleware.Middleware, database *database.Database, redis *Redis, log zerolog.Logger) {
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
					prefork := "Enabled"
					procs := runtime.GOMAXPROCS(0)
					if !cfg.App.Prefork {
						procs = 1
						prefork = "Disabled"
					}

					log.Info().Msgf("Version: %s", "-")
					log.Info().Msgf("Serve On: %s", cfg.ParseAddress())
					log.Info().Msgf("Prefork: %s", prefork)
					log.Info().Msgf("Handlers: %d", fiber.HandlersCount())
					log.Info().Msgf("Processes: %d", procs)
					log.Info().Msgf("PID: %d", os.Getpid())
				}

				// Listen the app (with TLS Support)
				if cfg.App.TLS.Enable {
					log.Info().Msg("TLS support was enabled.")

					if err := fiber.ListenTLS(cfg.App.Port, cfg.App.TLS.CertFile, cfg.App.TLS.KeyFile); err != nil {
						log.Error().Err(err).Msg("An unknown error occurred when to run server!")
					}
				}

				go func() {
					if err := fiber.Listen(":" + cfg.App.Port); err != nil {
						log.Error().Err(err).Msg("An unknown error occurred when to run server!")
					}
				}()

				redis.ConnectRedis()
				database.ConnectDatabase()

				//seeder := flag.Bool("seed", false, "seed the databases")
				//migrate := flag.Bool("migrate", false, "migrate the databases")
				//drop := flag.Bool("drop-all-tables", false, "drop all tables in the databases")
				//generateNecessaryData := flag.Bool("generate-necessary-data", false, "generating necessary data")
				//flag.Parse()
				//
				//if *migrate || *seeder || *drop || *generateNecessaryData {
				//	// read flag -migrate to migrate the database
				//	if *migrate {
				//		database.MigrateModels()
				//	}
				//	// read flag -generate-necessary-data to generate necessary data in the database
				//	if *generateNecessaryData {
				//		database.GenerateNecessaryData()
				//	}
				//	// read flag -seed to seed the database
				//	if *seeder {
				//		database.SeedModels()
				//	}
				//	// read flag -drop-all-tables to drop all tables in the database
				//	if *drop {
				//		database.DropTables()
				//	}
				//
				//	_ = fiber.Shutdown()
				//	database.ShutdownDatabase()
				//	os.Exit(0)
				//}

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
