package controller

import (
	"github.com/gofiber/fiber/v2"
	bService "go-fiber-starter/app/module/business/service"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/app/module/user/service"
	"go-fiber-starter/utils"
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
	BusinessUsers(c *fiber.Ctx) error
	InsertUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

func RestController(s service.IService, b bService.IService) IRestController {
	return &controller{s, b}
}

type controller struct {
	service  service.IService
	bService bService.IService
}

// Index
// @Summary      Get all users
// @Tags         Users
// @Security     Bearer
// @Router       /users [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Users
	req.Pagination = paginate
	req.Keyword = c.Query("Keyword")

	users, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: users,
		Meta: paging,
	})
}

// Show
// @Summary      Get one user
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "User ID"
// @Router       /users/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	user, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// Store
// @Summary      Create user
// @Tags         Users
// @Param 		 user body request.User true "User details"
// @Router       /users [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.User)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update
// @Summary      update user
// @Security     Bearer
// @Tags         Users
// @Param 		 user body request.User true "User details"
// @Param        id path int true "User ID"
// @Router       /users/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.User)
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
// @Summary      delete user
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "User ID"
// @Router       /users/:id [delete]
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

// BusinessUsers
// @Summary      Get one business users
// @Tags         Businesses
// @Security     Bearer
// @Param        id path int true "Business ID"
// @Router       /businesses/:businessID/users [get]
func (_i *controller) BusinessUsers(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrBadRequest
	}

	business, err := _i.bService.Show(businessID)
	if err != nil {
		return fiber.ErrNotFound
	}

	if user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.BusinessUsers
	req.Pagination = paginate
	req.BusinessID = businessID

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

	business, err := _i.bService.Show(businessID)
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

	business, err := _i.bService.Show(businessID)
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
