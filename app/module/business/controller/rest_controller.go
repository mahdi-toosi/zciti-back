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
	Users(c *fiber.Ctx) error
	InsertUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
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
// @Router       /businesses/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	business, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return c.JSON(business)
}

// Users
// @Summary      Get one business users
// @Tags         Businesses
// @Security     Bearer
// @Param        id path int true "Business ID"
// @Router       /businesses/:id/users [get]
func (_i *controller) Users(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Users
	req.BusinessID, err = utils.GetIntInParams(c, "id")
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.Pagination = paginate
	req.Keyword = c.Query("Keyword")

	users, paging, err := _i.service.Users(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: users,
		Meta: paging,
	})
}

// InsertUser
// @Summary      Insert one business user
// @Tags         Businesses
// @Security     Bearer
// @Param        businessId path int true "Business ID" ,userId path int true "User ID"
// @Router       /businesses/:businessID/users/:userID [post]
func (_i *controller) InsertUser(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := utils.GetIntInParams(c, "userID")
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

	if user.ID != userID && user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	err = _i.service.InsertUser(businessID, userID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"success"},
	})
}

// DeleteUser
// @Summary      Delete one business user
// @Tags         Businesses
// @Security     Bearer
// @Param        businessId path int true "Business ID" ,userId path int true "User ID"
// @Router       /businesses/:businessID/users/:userID [delete]
func (_i *controller) DeleteUser(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := utils.GetIntInParams(c, "userID")
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

	if user.ID != userID && user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	err = _i.service.DeleteUser(businessID, userID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"success"},
	})
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
// @Router       /businesses/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "id")
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
// @Router       /businesses/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "id")
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
