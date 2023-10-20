package internal

import (
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"go-fiber-starter/utils/config"
)

func NewMessageWay(cfg *config.Config) *MessageWay.App {
	return MessageWay.New(MessageWay.Config{
		ApiKey: cfg.Services.MessageWay.ApiKey,
	})
}
