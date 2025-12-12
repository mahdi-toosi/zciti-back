package internal

import (
	"go-fiber-starter/utils/config"

	baleBotApi "github.com/ghiac/bale-bot-api"
	"github.com/rs/zerolog/log"
)

type BaleBot struct {
	Bot          *baleBotApi.BotAPI
	LoggerChatID int64
	Connected    bool
}

func NewBaleBotLogger(cfg *config.Config) *BaleBot {
	if !cfg.App.Production {
		return &BaleBot{
			Bot:          nil,
			Connected:    false,
			LoggerChatID: cfg.Services.BaleBot.LoggerChatID,
		}
	}

	bot, err := baleBotApi.NewBaleBotAPI(cfg.Services.BaleBot.LoggerBotToken)
	if err != nil {
		log.Err(err).Msg("Error creating new BaleBot")
		return &BaleBot{
			Bot:          nil,
			Connected:    false,
			LoggerChatID: cfg.Services.BaleBot.LoggerChatID,
		}
	}

	bot.Debug = cfg.Services.BaleBot.Debug

	u := baleBotApi.NewUpdate(0)
	u.Timeout = 60
	// updates, err := bot.GetUpdatesChan(u)

	// if err != nil {
	//	return nil
	//}
	//for update := range updates {
	//	if update.Message == nil { // ignore any non-Message updates
	//		continue
	//	}
	//
	//	//if !update.Message.IsCommand() { // ignore any non-command Messages
	//	//	continue
	//	//}
	//
	//	msg := baleBotApi.NewMessage(update.Message.Chat.ID, "Hiii :)")
	//
	//	if _, err := bot.Send(msg); err != nil {
	//		log.Err(err).Msg("fail to send msg to bale bot")
	//	}
	//}

	return &BaleBot{
		Bot:          bot,
		Connected:    true,
		LoggerChatID: cfg.Services.BaleBot.LoggerChatID,
	}
}
