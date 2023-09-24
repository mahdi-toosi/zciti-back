package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/app/module/user/service"
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

// Index all Users
// @Summary      Get all users
// @Tags         Users
// @Security     Bearer
// @Router       /users [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.UsersRequest
	req.Pagination = paginate

	users, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User list successfully retrieved"},
		Data:     users,
		Meta:     paging,
	})
}

// Show one User
// @Summary      Get one user
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "User ID"
// @Router       /users/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	users, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully retrieved"},
		Data:     users,
	})
}

// Store user
// @Summary      Create user
// @Tags         Users
// @Param 		 user body request.UserRequest true "User details"
// @Router       /users [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.UserRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully created"},
	})
}

// Update user
// @Summary      update user
// @Security     Bearer
// @Tags         Users
// @Param 		 user body request.UserRequest true "User details"
// @Param        id path int true "User ID"
// @Router       /users/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	req := new(request.UserRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully updated"},
	})
}

// Delete user
// @Summary      delete user
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "User ID"
// @Router       /users/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully deleted"},
	})
}
