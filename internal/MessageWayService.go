package internal

import (
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/rs/zerolog"
	"go-fiber-starter/utils"
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
		utils.Log("sending message with this payload =>")
		utils.Log(req)
		//_s.logger.Warn().Msg("sending message with this payload => %v")
		//_s.logger.Info().Msgf("sending message with this payload => %v", req)

		return &MessageWay.SendResponse{
			Status:      "success",
			ReferenceID: "fakeReferenceID",
		}, nil
	}
	return _s.App.Send(req)
}

func NewMessageWay(cfg *config.Config, logger zerolog.Logger) *MessageWayService {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: cfg.Services.MessageWay.ApiKey,
	})

	return &MessageWayService{App: app, isProduction: cfg.App.Production, logger: logger}
}
