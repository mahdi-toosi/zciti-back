package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/config"
)

// Middleware is a struct that contains all the middleware functions
type Middleware struct {
	App *fiber.App
	Cfg *config.Config
}

func NewMiddleware(app *fiber.App, cfg *config.Config) *Middleware {
	return &Middleware{
		App: app,
		Cfg: cfg,
	}
}

// Register registers all the middleware functions
func (m *Middleware) Register() {
	// Add Extra Middlewares

	m.App.Use(helmet.New(helmet.Config{}))

	m.App.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Set("Access-Control-Allow-Origin", "*")
		c.Response().Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		c.Response().Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Response().Header.Set("Access-Control-Allow-Credentials", "true")
		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}
		// Continue to next middleware or handler
		return c.Next()
	})

	m.App.Use(limiter.New(limiter.Config{
		Next:       utils.IsEnabled(m.Cfg.Middleware.Limiter.Enable),
		Max:        m.Cfg.Middleware.Limiter.Max,
		Expiration: m.Cfg.Middleware.Limiter.Expiration * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return &fiber.Error{
				Code:    fiber.StatusTooManyRequests,
				Message: "تعداد درخواست شما به بیش از حد مجاز رسید.",
			}
		},
	}))

	m.App.Use(compress.New(compress.Config{
		Next:  utils.IsEnabled(m.Cfg.Middleware.Compress.Enable),
		Level: m.Cfg.Middleware.Compress.Level,
	}))

	m.App.Use(recover.New(recover.Config{
		Next: utils.IsEnabled(m.Cfg.Middleware.Recover.Enable),
	}))

	m.App.Use(pprof.New(pprof.Config{
		Next: utils.IsEnabled(m.Cfg.Middleware.Pprof.Enable),
	}))

	m.App.Use("/asset", filesystem.New(filesystem.Config{
		Next:   utils.IsEnabled(m.Cfg.Middleware.FileSystem.Enable),
		Root:   http.Dir(m.Cfg.Middleware.FileSystem.Root),
		Browse: m.Cfg.Middleware.FileSystem.Browse,
		MaxAge: m.Cfg.Middleware.FileSystem.MaxAge,
	}))

	m.App.Get(m.Cfg.Middleware.Monitor.Path, monitor.New(monitor.Config{
		Next: utils.IsEnabled(m.Cfg.Middleware.Monitor.Enable),
	}))
}
