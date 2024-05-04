package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/business/request"
	res "go-fiber-starter/app/module/business/response"
	"go-fiber-starter/app/module/business/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Types(c *fiber.Ctx) error
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

// Index
// @Summary      Get all businesses
// @Tags         Businesses
// @Security     Bearer
// @Router       /businesses [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	// TODO we need new endpoint for hiding business phone numbers
	// TODO only admins can access this endpoint

	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Businesses
	req.Pagination = paginate
	req.Keyword = c.Query("Keyword")

	businesses, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: businesses,
		Meta: paging,
	})
}

// Types
// @Summary      Get businesses types
// @Tags         Businesses
// @Security     Bearer
// @Router       /businesses/types [get]
func (_i *controller) Types(c *fiber.Ctx) error {
	var types []res.BusinessTypesOption //nolint:prealloc
	for value, label := range schema.TypeDisplayProxy {
		types = append(types, res.BusinessTypesOption{
			Label: label,
			Value: value,
		})
	}

	return c.JSON(types)
}

// Show
// @Summary      Get one business
// @Tags         Businesses
// @Security     Bearer
// @Param        id path int true "Business ID"
// @Router       /businesses/:businessID [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	// TODO remove it
	if businessID == 1 {
		return fiber.ErrForbidden
	}

	business, err := _i.service.Show(businessID)
	if err != nil {
		return err
	}

	return c.JSON(business)
}

// Store
// @Summary      Create business
// @Tags         Businesses
// @Param 		 business body request.Business true "Business details"
// @Router       /businesses [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.Business)
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
// @Summary      update business
// @Security     Bearer
// @Tags         Businesses
// @Param 		 business body request.Business true "Business details"
// @Param        id path int true "Business ID"
// @Router       /businesses/:businessID [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrBadRequest
	}

	business, err := _i.service.Show(businessID)
	if err != nil {
		return fiber.ErrNotFound
	}

	if user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	req := new(request.Business)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(businessID, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete
// @Summary      delete business
// @Tags         Businesses
// @Security     Bearer
// @Param        id path int true "Business ID"
// @Router       /businesses/:businessID [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrBadRequest
	}

	business, err := _i.service.Show(businessID)
	if err != nil {
		return fiber.ErrNotFound
	}

	if user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	err = _i.service.Destroy(businessID)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
