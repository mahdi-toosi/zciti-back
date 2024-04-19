package controller

import (
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/notification/request"
	"go-fiber-starter/app/module/notification/service"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
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

// Index all Notifications
// @Summary      Get all notifications
// @Tags         Notifications
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/notifications [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Notifications
	req.BusinessID = businessID
	req.Pagination = paginate

	notifications, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: notifications,
		Meta: paging,
	})
}

// Show one Notification
// @Summary      Get one notification
// @Tags         Notifications
// @Security     Bearer
// @Param        id path int true "Notification ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/notifications/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	notification, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(notification)
}

// Store notification
// @Summary      Create notification
// @Tags         Notifications
// @Param 		 notification body request.Notification true "Notification details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/notifications [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.Notification)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.BusinessID = businessID
	err = _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update notification
// @Summary      update notification
// @Security     Bearer
// @Tags         Notifications
// @Param 		 notification body request.Notification true "Notification details"
// @Param        id path int true "Notification ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/notifications/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Notification)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.BusinessID = businessID
	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete notification
// @Summary      delete notification
// @Tags         Notifications
// @Security     Bearer
// @Param        id path int true "Notification ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/notifications/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
