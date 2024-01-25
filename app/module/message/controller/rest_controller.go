package controller

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	businessResponse "go-fiber-starter/app/module/business/response"
	businessService "go-fiber-starter/app/module/business/service"
	"go-fiber-starter/app/module/message/request"
	msgResponse "go-fiber-starter/app/module/message/response"
	"go-fiber-starter/app/module/message/service"
	msgRoomResponse "go-fiber-starter/app/module/messageRoom/response"
	messageRoomService "go-fiber-starter/app/module/messageRoom/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"time"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Stream(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService, ms messageRoomService.IService, bs businessService.IService) IRestController {
	return &controller{service: s, messageRoomService: ms, businessService: bs}
}

type controller struct {
	service            service.IService
	businessService    businessService.IService
	messageRoomService messageRoomService.IService
}

// Index
// @Summary      Get all messages in room
// @Tags         Messages
// @Security     Bearer
// @Router       /messages/:businessID [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	token := c.Query("Token", "")

	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrForbidden
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Messages
	req.UserID = user.ID
	req.Pagination = paginate
	req.BusinessID = businessID

	var business *schema.Business
	var paging paginator.Pagination
	var msgRoom *schema.MessageRoom
	var messages []*msgResponse.Message

	if tokenRoomData, err := _i.messageRoomService.IsTokenValid(token); err != nil {
		messages, paging, msgRoom, business, err = _i.service.Index(req, nil)
		if err != nil {
			return fiber.ErrBadRequest
		}

		token, err = _i.messageRoomService.GenerateToken(msgRoomResponse.FromDomain(msgRoom, nil), businessResponse.FromDomain(business))
		if err != nil {
			return err
		}
	} else {
		msgRoom = &schema.MessageRoom{
			ID:         tokenRoomData.ID,
			UserID:     tokenRoomData.UserID,
			Status:     tokenRoomData.Status,
			BusinessID: tokenRoomData.BusinessID,
		}
		messages, paging, _, _, err = _i.service.Index(req, &msgRoom.ID)
		if err != nil {
			return fiber.ErrBadRequest
		}
	}

	if messages == nil {
		messages = []*msgResponse.Message{}
	}
	return response.Resp(c, response.Response{
		Data: response.Response{
			Data: messages,
			Meta: paging,
		},
		Meta: msgRoomResponse.FromDomain(msgRoom, &token),
	})
}

func (_i *controller) Stream(c *fiber.Ctx) error {
	token := c.Query("Token", "")
	if token == "" {
		return fiber.ErrBadRequest
	}

	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrForbidden
	}

	var req request.Messages
	req.UserID = user.ID
	req.BusinessID = businessID

	var messages []*msgResponse.Message

	tokenRoomData, err := _i.messageRoomService.IsTokenValid(token)
	if err != nil {
		if err != nil {
			return fiber.ErrBadRequest
		}
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		var i int
		for {
			i++
			msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
			_, err := fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			if err != nil {
				continue
			}

			err = w.Flush()
			if err != nil {
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	})

	for {
		messages = _i.service.CheckForNewMessages(tokenRoomData.ID)
		log.Debug().Msgf("%+v", messages)

		if messages == nil {
			messages = []*msgResponse.Message{}
		}
		//return response.Resp(c, response.Response{
		//	Data: messages,
		//})
		time.Sleep(2 * time.Second)
	}
	return nil

	//ID:         tokenRoomData.ID,
	//UserID:     tokenRoomData.UserID,
	//Status:     tokenRoomData.Status,
	//BusinessID: tokenRoomData.BusinessID,
}

// Store
// @Summary      Create message
// @Tags         Messages
// @Param 		 message body request.Message true "Message details"
// @Router       /messages [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.Message)

	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	token := req.Token

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrForbidden
	}

	if tokenRoomData, err := _i.messageRoomService.IsTokenValid(token); err != nil {
		msgRoom, err := _i.messageRoomService.ShowByID(req.RoomID)
		if err != nil {
			return fiber.ErrBadRequest
		}

		business, err := _i.businessService.Show(msgRoom.BusinessID)
		if err != nil {
			return fiber.ErrBadRequest
		}

		token, err = _i.messageRoomService.GenerateToken(msgRoomResponse.FromDomain(msgRoom, nil), business)
		if err != nil {
			return err
		}
	} else if !tokenRoomData.HasMember(user.ID) {
		return fiber.ErrBadRequest
	}

	msg, err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: msgResponse.StoreMessage{Token: token, ID: msg.ID},
	})
}

// Update
// @Summary      update message
// @Security     Bearer
// @Tags         Messages
// @Param 		 message body request.Message true "Message details"
// @Param        id path int true "Message ID"
// @Router       /messages/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Message)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete
// @Summary      delete message
// @Tags         Messages
// @Security     Bearer
// @Param        id path int true "Message ID"
// @Router       /messages/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
