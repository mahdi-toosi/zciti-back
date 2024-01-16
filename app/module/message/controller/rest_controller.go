package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/message/request"
	msgResponse "go-fiber-starter/app/module/message/response"
	"go-fiber-starter/app/module/message/service"
	msgRoomResponse "go-fiber-starter/app/module/messageRoom/response"
	messageRoomService "go-fiber-starter/app/module/messageRoom/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService, ms messageRoomService.IService) IRestController {
	return &controller{service: s, messageRoomService: ms}
}

type controller struct {
	service            service.IService
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
			return err
		}

		token, err = _i.messageRoomService.GenerateToken(msgRoom, business)
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
			return err
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

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
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
