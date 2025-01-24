package internal

import (
	baleBotApi "github.com/ghiac/bale-bot-api"
	"go-fiber-starter/utils/config"
)

type BaleBot struct {
	Bot          *baleBotApi.BotAPI
	LoggerChatID int64
}

func NewBaleBotLogger(cfg *config.Config) *BaleBot {
	bot, err := baleBotApi.NewBaleBotAPI(cfg.Services.BaleBot.LoggerBotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = cfg.Services.BaleBot.Debug

	u := baleBotApi.NewUpdate(0)
	u.Timeout = 60
	// updates, err := bot.GetUpdatesChan(u)

	// if err != nil {
	//	return nil
	//}
	//for update := range updates {
	//	utils.Log(update)
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
		LoggerChatID: cfg.Services.BaleBot.LoggerChatID,
	}
}
