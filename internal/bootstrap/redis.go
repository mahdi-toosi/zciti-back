package bootstrap

import (
	"github.com/gofiber/storage/redis/v3"
	"github.com/rs/zerolog"
	"go-fiber-starter/utils/config"
)

type Redis struct {
	Storage *redis.Storage
	Log     zerolog.Logger
	Cfg     *config.Config
}

func NewRedis(cfg *config.Config, log zerolog.Logger) *Redis {
	return &Redis{
		Cfg: cfg,
		Log: log,
	}
}

func (_db *Redis) ConnectRedis() {
	storage := redis.New(redis.Config{
		URL: _db.Cfg.Redis.Url,
	})

	_db.Log.Info().Msg("Connected the Redis successfully!")

	_db.Storage = storage
}
