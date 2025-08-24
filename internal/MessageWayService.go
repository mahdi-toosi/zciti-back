package internal

import (
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"go-fiber-starter/utils/config"
)

type MessageWayService struct {
	isProduction bool
	*MessageWay.App
}

// Send Overriding the MessageWay Send method
func (_s *MessageWayService) Send(req MessageWay.Message) (*MessageWay.SendResponse, error) {
	if !_s.isProduction {
		req.Mobile = "09380338494"
	}
	return _s.App.Send(req)
}

func NewMessageWay(cfg *config.Config) *MessageWayService {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: cfg.Services.MessageWay.ApiKey,
	})
	return &MessageWayService{App: app, isProduction: cfg.App.Production}
}
