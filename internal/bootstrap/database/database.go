package database

import (
	"github.com/rs/zerolog"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/database/seeds"
	"go-fiber-starter/utils/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB  *gorm.DB
	Log zerolog.Logger
	Cfg *config.Config
}

func NewDatabase(cfg *config.Config, log zerolog.Logger) *Database {
	db := &Database{
		Cfg: cfg,
		Log: log,
	}

	return db
}

func (_db *Database) ConnectDatabase() {
	conn, err := gorm.Open(postgres.Open(_db.Cfg.DB.Postgres.DSN), &gorm.Config{})
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to connect the database!")
	} else {
		_db.Log.Info().Msg("Connected the database successfully!")
	}

	_db.DB = conn
}

func (_db *Database) ShutdownDatabase() {
	sqlDB, err := _db.DB.DB()
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to shutdown the database!")
	} else {
		_db.Log.Info().Msg("Shutdown the database successfully!")
	}
	_err := sqlDB.Close()
	if _err != nil {
		return
	}
}

func (_db *Database) MigrateModels() {
	if err := _db.DB.AutoMigrate(schema.Models()...); err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to migrate the database!")
		panic("An unknown error occurred when to migrate the database!")
	}
}

func (_db *Database) SeedModels() {
	for _, model := range seeds.Seeders() {
		count, err := model.Count(_db.DB)
		if err != nil {
			_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the database!")
		}

		if count == 0 {
			if err := model.Seed(_db.DB); err != nil {
				_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the database!")
			}

			_db.Log.Info().Msg("Seeded the database successfully!")
		} else {
			_db.Log.Info().Msg("Database is already seeded!")
		}
	}
}

func (_db *Database) DropTables() {
	for _, model := range schema.Models() {
		if err := _db.DB.Migrator().DropTable(model); err != nil {
			_db.Log.Error().Err(err).Msg("An unknown error occurred when to drop table in the database!")
		}
	}

	schema.DropExtraCommands(_db.DB)
}
