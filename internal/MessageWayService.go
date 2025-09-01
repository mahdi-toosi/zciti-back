package internal

import (
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/rs/zerolog"
	"go-fiber-starter/utils/config"
)

type MessageWayService struct {
	isProduction bool
	*MessageWay.App
	logger zerolog.Logger
}

// Send Overriding the MessageWay Send method
func (_s *MessageWayService) Send(req MessageWay.Message) (*MessageWay.SendResponse, error) {
	if !_s.isProduction {
		// sms mode
		//req.Mobile = "09380338494"
		// logger mode
		_s.logger.Warn().Msgf("sending message with this payload => %v", req)
		return nil, nil
	}
	return _s.App.Send(req)
}

func NewMessageWay(cfg *config.Config, logger zerolog.Logger) *MessageWayService {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: cfg.Services.MessageWay.ApiKey,
	})
	return &MessageWayService{App: app, isProduction: cfg.App.Production, logger: logger}
}
