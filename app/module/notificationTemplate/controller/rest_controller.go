package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/notificationTemplate/request"
	"go-fiber-starter/app/module/notificationTemplate/service"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"strconv"
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

// Index all NotificationTemplates
// @Summary      Get all notificationTemplates
// @Tags         NotificationTemplates
// @Security     Bearer
// @Router       /notificationTemplates [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Index
	req.Pagination = paginate

	notificationTemplates, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: notificationTemplates,
		Meta: paging,
	})
}

// Show one NotificationTemplate
// @Summary      Get one NotificationTemplate
// @Tags         NotificationTemplates
// @Security     Bearer
// @Param        id path int true "NotificationTemplate ID"
// @Router       /notificationTemplates/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	notificationTemplate, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return c.JSON(notificationTemplate)
}

// Store NotificationTemplate
// @Summary      Create NotificationTemplate
// @Tags         NotificationTemplates
// @Param 		 NotificationTemplate body request.NotificationTemplate true "NotificationTemplate details"
// @Router       /notificationTemplates [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.NotificationTemplate)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update NotificationTemplate
// @Summary      update NotificationTemplate
// @Security     Bearer
// @Tags         NotificationTemplates
// @Param 		 NotificationTemplate body request.NotificationTemplate true "NotificationTemplate details"
// @Param        id path int true "NotificationTemplate ID"
// @Router       /notificationTemplates/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	req := new(request.NotificationTemplate)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete NotificationTemplate
// @Summary      delete NotificationTemplate
// @Tags         NotificationTemplates
// @Security     Bearer
// @Param        id path int true "NotificationTemplate ID"
// @Router       /notificationTemplates/:id [delete]
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
