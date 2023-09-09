package database

import (
	"github.com/bangadam/go-fiber-starter/app/database/schema"
	"github.com/bangadam/go-fiber-starter/app/database/seeds"
	"github.com/bangadam/go-fiber-starter/utils/config"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setup database with gorm
type Database struct {
	DB  *gorm.DB
	Log zerolog.Logger
	Cfg *config.Config
}

type Seeder interface {
	Seed(*gorm.DB) error
	Count(*gorm.DB) (int, error)
}

func NewDatabase(cfg *config.Config, log zerolog.Logger) *Database {
	db := &Database{
		Cfg: cfg,
		Log: log,
	}

	return db
}

// connect database
func (_db *Database) ConnectDatabase() {
	conn, err := gorm.Open(postgres.Open(_db.Cfg.DB.Postgres.DSN), &gorm.Config{})
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to connect the database!")
	} else {
		_db.Log.Info().Msg("Connected the database successfully!")
	}

	_db.DB = conn
}

// shutdown database
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

// migrate models
func (_db *Database) MigrateModels() {
	if err := _db.DB.AutoMigrate(
		Models()...,
	); err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to migrate the database!")
	}
}

// list of models for migration
func Models() []interface{} {
	return []interface{}{
		schema.User{},
	}
}

// seed data
func (_db *Database) SeedModels() {
	seeders := []Seeder{
		seeds.UserSeeder{},
	}

	for _, seed := range seeders {
		count, err := seed.Count(_db.DB)
		if err != nil {
			_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the database!")
		}

		if count == 0 {
			if err := seed.Seed(_db.DB); err != nil {
				_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the database!")
			}

			_db.Log.Info().Msg("Seeded the database successfully!")
		} else {
			_db.Log.Info().Msg("Database is already seeded!")
		}
	}

	_db.Log.Info().Msg("Seeded the database successfully!")
}
