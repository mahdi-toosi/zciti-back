package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// app struct config
type app = struct {
	Name           string        `toml:"name"`
	FrontendDomain string        `toml:"frontendDomain"`
	BackendDomain  string        `toml:"backendDomain"`
	Port           string        `toml:"port"`
	PrintRoutes    bool          `toml:"print-routes"`
	Prefork        bool          `toml:"prefork"`
	Production     bool          `toml:"production"`
	IdleTimeout    time.Duration `toml:"idle-timeout"`
	TLS            struct {
		Enable   bool   `toml:"enable"`
		CertFile string `toml:"cert-file"`
		KeyFile  string `toml:"key-file"`
	}
}

// db struct config
type db = struct {
	Main struct {
		Url string `toml:"url"`
	}
	Chat struct {
		Url string `toml:"url"`
	}
}

// redis struct config
type redis = struct {
	Url string `toml:"url"`
}

// log struct config
type logger = struct {
	TimeFormat string        `toml:"time-format"`
	Level      zerolog.Level `toml:"level"`
	Prettier   bool          `toml:"prettier"`
}

type services = struct {
	MessageWay struct {
		ApiKey string `toml:"apiKey"`
	}

	ZarinPal struct {
		MerchantID string `toml:"merchantId"`
		Sandbox    bool   `toml:"sandbox"`
	}

	GoogleRecaptcha struct {
		SecretKey string `toml:"secretKey"`
	}

	BaleBot struct {
		Debug          bool   `toml:"debug"`
		LoggerChatID   int64  `toml:"loggerChatID"`
		LoggerBotToken string `toml:"loggerBotToken"`
	}
}

// middleware
type middleware = struct {
	Compress struct {
		Enable bool
		Level  compress.Level
	}

	Recover struct {
		Enable bool
	}

	Monitor struct {
		Enable bool
		Path   string
	}

	Pprof struct {
		Enable bool
	}

	Limiter struct {
		Enable     bool
		Max        int
		Expiration time.Duration `toml:"expiration_seconds"`
	}

	FileSystem struct {
		Enable      bool
		Browse      bool
		MaxAge      int `toml:"max_age"`
		Index       string
		Root        string
		PrivateRoot string `toml:"private_root"`
	}

	Jwt Jwt
}

type Jwt struct {
	Secret     string        `toml:"secret"`
	Expiration time.Duration `toml:"expiration_seconds"`
}

type Config struct {
	App        app
	DB         db
	Redis      redis
	Logger     logger
	Middleware middleware
	Services   services
}

// ParseConfig func to parse config
func ParseConfig(name string, debug ...bool) (*Config, error) {
	var (
		contents *Config
		file     []byte
		err      error
	)

	if len(debug) > 0 {
		file, err = os.ReadFile(name)
	} else {
		_, b, _, _ := runtime.Caller(0)
		// get base path
		path := filepath.Dir(filepath.Dir(filepath.Dir(b)))
		file, err = os.ReadFile(filepath.Join(path, "./config/", name+".toml"))
	}

	if err != nil {
		return &Config{}, err
	}

	err = toml.Unmarshal(file, &contents)

	return contents, err
}

// NewConfig initialize config
func NewConfig() *Config {
	config, err := ParseConfig("zciti")
	if err != nil && !fiber.IsChild() {
		// panic if config is not found
		log.Panic().Err(err).Msg("config not found")
	}

	if config.Middleware.Jwt.Secret == "" {
		panic("JWT secret is not set")
	}

	return config
}

// ParseAddress func to parse address
func (c Config) ParseAddress() string {
	return c.App.BackendDomain
}
