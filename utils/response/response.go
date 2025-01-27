package response

import (
	"fmt"
	baleBotApi "github.com/ghiac/bale-bot-api"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Alias for any slice
type Messages = []any

type Error struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

// error makes it compatible with the error interface
func (e *Error) Error() string {
	return fmt.Sprint(e.Message)
}

// A struct to return normal response
type Response struct {
	Code     int      `json:",omitempty"`
	Messages Messages `json:",omitempty"`
	Data     any      `json:",omitempty"`
	Meta     any      `json:",omitempty"`
	URL      string   `json:",omitempty"`
	Token    []string `json:",omitempty"`
}

var IsProduction bool

// Default error handler
var ErrorHandler = func(c *fiber.Ctx, err error, baleBot *internal.BaleBot) error {
	resp := Response{
		Code: fiber.StatusInternalServerError,
	}

	//_, ok := err.(validator.ValidationErrors)

	// handle errors
	if c, ok := err.(validator.ValidationErrors); ok {
		resp.Code = fiber.StatusUnprocessableEntity
		resp.Messages = Messages{removeTopStruct(c.Translate(trans))}
	} else if c, ok := err.(*fiber.Error); ok {
		resp.Code = c.Code
		resp.Messages = Messages{c.Message}
	} else if c, ok := err.(*Error); ok {
		resp.Code = c.Code
		resp.Messages = Messages{c.Message}

		// for ent and other errors
		if resp.Messages == nil {
			resp.Messages = Messages{err}
		}
	} else {
		resp.Messages = Messages{err.Error()}
	}
	if !IsProduction {
		log.Error().Err(err).Msg("From: Fiber's error handler")
	}

	if baleBot.Connected {
		baleBotMsgPayload := resp
		baleBotMsgPayload.URL = c.Request().URI().String()
		baleBotMsgPayload.Token = c.GetReqHeaders()["Authorization"]
		msg := baleBotApi.NewMessage(baleBot.LoggerChatID, utils.PrettyJSON(baleBotMsgPayload))
		if _, err := baleBot.Bot.Send(msg); err != nil {
			log.Err(err).Msg("fail to send msg to bale bot")
		}
	}

	return Resp(c, resp)
}

// function to return pretty json response
func Resp(c *fiber.Ctx, resp Response) error {
	// set Data
	if resp.Data == nil {
		resp.Data = []any{}
	}
	// set status
	if resp.Code == 0 {
		resp.Code = fiber.StatusOK
	}
	c.Status(resp.Code)
	// return response
	return c.JSON(resp)
}

// remove unnecessary fields from validator message
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}

	for field, msg := range fields {
		stripStruct := field[strings.Index(field, ".")+1:]
		res[stripStruct] = msg
	}

	return res
}
