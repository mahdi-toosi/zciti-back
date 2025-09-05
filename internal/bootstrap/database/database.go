package database //nolint:typecheck

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/database/seeds"
	"go-fiber-starter/utils/config"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Main *gorm.DB
	// Chat *gorm.DB
	Log zerolog.Logger
	Cfg *config.Config
}

func NewDatabase(cfg *config.Config, log zerolog.Logger) *Database {
	return &Database{
		Cfg: cfg,
		Log: log,
	}
}

func (_db *Database) ConnectDatabase() {
	mainDB, err := gorm.Open(postgres.Open(_db.Cfg.DB.Main.Url), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		QueryFields:            _db.Cfg.App.Production,
	})
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to connect the *Main* database!")
	} else {
		_db.Log.Info().Msg("Connected the *Main* database successfully!")
	}

	_db.Main = mainDB

	//chatDB, err := gorm.Open(postgres.Open(_db.Cfg.DB.Chat.Url), &gorm.Config{})
	//if err != nil {
	//	_db.Log.Error().Err(err).Msg("An unknown error occurred when to connect the *Chat* database!")
	//} else {
	//	_db.Log.Info().Msg("Connected the *Chat* database successfully!")
	//}
	//
	//_db.Chat = chatDB
}

func (_db *Database) ShutdownDatabase() {
	mainDB, err := _db.Main.DB()
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to shutdown the *Main* database!")
	} else {
		_db.Log.Info().Msg("Shutdown the *Main* database successfully!")
	}
	_err := mainDB.Close()
	if _err != nil {
		return
	}

	//chatDB, err := _db.Chat.DB()
	//if err != nil {
	//	_db.Log.Error().Err(err).Msg("An unknown error occurred when to shutdown the *Chat* database!")
	//} else {
	//	_db.Log.Info().Msg("Shutdown the *Chat* database successfully!")
	//}
	//_err = chatDB.Close()
	//if _err != nil {
	//	return
	//}
}

func (_db *Database) setUUIDExtension() {
	_ = _db.Main.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
}

func (_db *Database) MigrateModels() {
	_db.setUUIDExtension()
	if err := _db.Main.AutoMigrate(schema.MainDBModels()...); err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when to migrate the *Main* database!")
		panic("An unknown error occurred when to migrate the *Main* database!")
	}

	//if err := _db.Chat.AutoMigrate(schema.ChatDBModels()...); err != nil {
	//	_db.Log.Error().Err(err).Msg("An unknown error occurred when to migrate the *Chat* database!")
	//	panic("An unknown error occurred when to migrate the *Chat* database!")
	//}
}

func (_db *Database) GenerateNecessaryData() {
	err := seeds.GenerateNecessaryData(_db.Main)
	if err != nil {
		_db.Log.Error().Err(err).Msg("An unknown error occurred when generating necessary data wa running!")
	}
}

func (_db *Database) SeedModels() {
	for _, model := range seeds.MainDBSeeders() {
		count, err := model.Count(_db.Main)
		if err != nil {
			_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the *Main* database!")
		}

		if count == 0 {
			if err := model.Seed(_db.Main); err != nil {
				_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the *Main* database!")
			}
		}
	}

	_db.Log.Info().Msg("*Main* Seeded the database successfully!")

	//for _, model := range seeds.ChatDBSeeders() {
	//	count, err := model.Count(_db.Chat)
	//	if err != nil {
	//		_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the *Chat* database!")
	//	}
	//
	//	if count == 0 {
	//		if err := model.Seed(_db.Chat); err != nil {
	//			_db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the *Chat* database!")
	//		}
	//	}
	//}
	//
	//_db.Log.Info().Msg("*Chat* Seeded the database successfully!")
}

func (_db *Database) DropTables() {
	for _, model := range schema.MainDBModels() {
		if err := _db.Main.Migrator().DropTable(model); err != nil {
			_db.Log.Error().Err(err).Msg("An unknown error occurred when to drop table in the *Main* database!")
		}
	}

	schema.MainDBDropExtraCommands(_db.Main)

	//for _, model := range schema.ChatDBModels() {
	//	if err := _db.Chat.Migrator().DropTable(model); err != nil {
	//		_db.Log.Error().Err(err).Msg("An unknown error occurred when to drop table in the *Chat* database!")
	//	}
	//}
	//
	//schema.ChatDBDropExtraCommands(_db.Chat)
}
