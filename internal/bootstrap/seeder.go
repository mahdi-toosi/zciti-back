package bootstrap

import (
	"context"
	"flag"
	"github.com/rs/zerolog"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/config"
	"go.uber.org/fx"
	"os"
)

func Seeder(lifecycle fx.Lifecycle, cfg *config.Config, database *database.Database, log zerolog.Logger) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Custom Startup Messages
				database.ConnectDatabase()

				seeder := flag.Bool("seed", false, "seed the databases")
				migrate := flag.Bool("migrate", false, "migrate the databases")
				drop := flag.Bool("drop-all-tables", false, "drop all tables in the databases")
				deleteFilesInStorage := flag.Bool("delete-files-in-storage", false, "generating necessary data")
				generateNecessaryData := flag.Bool("generate-necessary-data", false, "generating necessary data")
				flag.Parse()

				if *migrate || *seeder || *drop || *generateNecessaryData || *deleteFilesInStorage {
					// read flag -drop-all-tables to drop all tables in the database
					if *drop {
						database.DropTables()
					}
					// read flag -migrate to migrate the database
					if *migrate {
						database.MigrateModels()
					}
					// read flag -generate-necessary-data to generate necessary data in the database
					if *generateNecessaryData {
						database.GenerateNecessaryData()
					}
					// read flag -seed to seed the database
					if *seeder {
						database.SeedModels()
					}

					if *deleteFilesInStorage {
						if err := utils.DeleteFoldersInDirectory(cfg.Middleware.FileSystem.Root); err != nil {
							return err
						}
						if err := utils.DeleteFoldersInDirectory(cfg.Middleware.FileSystem.PrivateRoot); err != nil {
							return err
						}
					}

					database.ShutdownDatabase()
					os.Exit(0)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
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
