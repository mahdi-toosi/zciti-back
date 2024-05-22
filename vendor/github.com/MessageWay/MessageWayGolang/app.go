package MessageWay

import (
	"net/http"
	"time"
)

const ApiBaseUrl = "https://api.msgway.com"

func New(c Config) *App {
	if c.ApiKey == "" {
		panic("apiKey not set")
	}
	app := &App{config: c}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 5
	t.MaxConnsPerHost = 5
	t.MaxIdleConnsPerHost = 5

	app.client = &http.Client{
		Timeout:   5 * time.Second,
		Transport: t,
	}
	return app
}
