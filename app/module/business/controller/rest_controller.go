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
	OperatorShow(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	MyBusinesses(c *fiber.Ctx) error
	MenuItems(c *fiber.Ctx) error
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

	user, _ := utils.GetAuthenticatedUser(c)

	// TODO remove it
	if businessID == 1 && !user.IsAdmin() {
		return fiber.ErrForbidden
	}
	business, err := _i.service.Show(businessID, schema.URUser)
	if err != nil {
		return err
	}

	return c.JSON(business)
}

// OperatorShow
// @Summary      Get operator business
// @Tags         Businesses
// @Security     Bearer
// @Param        id path int true "Business ID"
// @Router       operator/businesses/:businessID [get]
func (_i *controller) OperatorShow(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	user, _ := utils.GetAuthenticatedUser(c)

	// TODO remove it
	if businessID == 1 && !user.IsAdmin() {
		return fiber.ErrForbidden
	}
	business, err := _i.service.Show(businessID, schema.URBusinessOwner)
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
		return err
	}

	business, err := _i.service.Show(businessID, schema.URUser)
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
		return err
	}

	business, err := _i.service.Show(businessID, schema.URBusinessOwner)
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

// MyBusinesses
// @Summary      Get user businesses
// @Tags         Businesses
// @Security     Bearer
// @Router       /user/businesses [get]
func (_i *controller) MyBusinesses(c *fiber.Ctx) error {
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if len(user.Permissions) == 0 {
		return c.JSON([]string{})
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Businesses
	paginate.Page = 0
	req.Pagination = paginate
	for ID := range user.Permissions {
		req.IDs = append(req.IDs, ID)
	}

	businesses, _, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return c.JSON(businesses)
}

// MenuItems
// @Summary      get user business menu items
// @Security     Bearer
// @Tags         Businesses
// @Param        id path int true "Business ID"
// @Router       /user/businesses/:businessID/menu-items [get]
func (_i *controller) MenuItems(c *fiber.Ctx) error {
	user, _ := utils.GetAuthenticatedUser(c)

	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	menuItems, err := _i.service.RoleMenuItems(businessID, user)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: menuItems,
	})
}
