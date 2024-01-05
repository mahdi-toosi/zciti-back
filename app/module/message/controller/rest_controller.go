package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/message/request"
	"go-fiber-starter/app/module/message/service"
	msgRoomResponse "go-fiber-starter/app/module/messageRoom/response"
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

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index
// @Summary      Get all messages in room
// @Tags         Messages
// @Security     Bearer
// @Router       /messages/:businessID [get]
func (_i *controller) Index(c *fiber.Ctx) error {
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
	req.Pagination = paginate
	req.UserID = user.ID
	req.BusinessID = businessID

	messages, paging, msgRoom, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: response.Response{
			Data: messages,
			Meta: paging,
		},
		Meta: msgRoomResponse.FromDomain(msgRoom),
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
