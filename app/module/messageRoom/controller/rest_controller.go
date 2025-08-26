package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/module/messageRoom/request"
	"go-fiber-starter/app/module/messageRoom/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index
// @Summary      Get all messageRooms
// @Tags         MessageRooms
// @Security     Bearer
// @Router       /message-rooms [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	userID, err := utils.GetUintInQueries(c, "userID")
	if err != nil {
		return err
	}
	businessID, err := utils.GetUintInQueries(c, "businessID")
	if err != nil {
		return err
	}

	log.Debug().Msgf("%+v", userID, businessID)

	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.MessageRooms
	req.Pagination = paginate

	messageRooms, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: messageRooms,
		Meta: paging,
	})
}

// Delete
// @Summary      delete messageRoom
// @Tags         MessageRooms
// @Security     Bearer
// @Param        id path int true "MessageRoom ID"
// @Router       /message-rooms/:id [delete]
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
