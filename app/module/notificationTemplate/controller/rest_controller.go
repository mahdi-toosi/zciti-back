package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/notificationTemplate/request"
	"go-fiber-starter/app/module/notificationTemplate/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"strconv"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Keywords(c *fiber.Ctx) error
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
// @Summary      Get all notificationTemplates
// @Tags         NotificationTemplates
// @Security     Bearer
// @Router       /business/:businessID/notificationTemplates [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Index
	req.Pagination = paginate
	req.BusinessID = businessID

	notificationTemplates, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: notificationTemplates,
		Meta: paging,
	})
}

// Keywords
// @Summary      Get NotificationTemplate keywords
// @Tags         NotificationTemplates
// @Security     Bearer
// @Router       /business/:businessID/notificationTemplates/keywords [get]
func (_i *controller) Keywords(c *fiber.Ctx) error {
	var keywords = map[string]string{
		"FirstName":     "نام کاربر",
		"LastName":      "نام خانوادگی کاربر",
		"FullName":      "نام و نام خانوادگی کاربر",
		"BusinessName":  "نام کسب و کار",
		"BusinessPhone": "تلفن کسب و کار",
	}

	return c.JSON(keywords)
}

// Store
// @Summary      Create NotificationTemplate
// @Tags         NotificationTemplates
// @Param 		 NotificationTemplate body request.NotificationTemplate true "NotificationTemplate details"
// @Router       /business/:businessID/notificationTemplates [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	req := new(request.NotificationTemplate)
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

// Update
// @Summary      update NotificationTemplate
// @Security     Bearer
// @Tags         NotificationTemplates
// @Param 		 NotificationTemplate body request.NotificationTemplate true "NotificationTemplate details"
// @Param        id path int true "NotificationTemplate ID"
// @Router       /business/:businessID/notificationTemplates/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	req := new(request.NotificationTemplate)
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

// Delete
// @Summary      delete NotificationTemplate
// @Tags         NotificationTemplates
// @Security     Bearer
// @Param        id path int true "NotificationTemplate ID"
// @Router       /business/:businessID/notificationTemplates/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.service.Destroy(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
