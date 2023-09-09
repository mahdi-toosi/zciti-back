package controller

import (
	"strconv"

	"github.com/bangadam/go-fiber-starter/app/module/user/request"
	"github.com/bangadam/go-fiber-starter/app/module/user/service"
	"github.com/bangadam/go-fiber-starter/utils/paginator"
	"github.com/bangadam/go-fiber-starter/utils/response"
	"github.com/gofiber/fiber/v2"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(userService service.IService) IRestController {
	return &controller{
		userService,
	}
}

type controller struct {
	userService service.IService
}

// Index all Users
// @Summary      Get all users
// @Description  API for getting all users
// @Tags         Task
// @Security     Bearer
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.UsersRequest
	req.Pagination = paginate

	users, paging, err := _i.userService.All(req)
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
// @Description  API for getting one user
// @Tags         Task
// @Security     Bearer
// @Param        id path int true "User ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	users, err := _i.userService.Show(id)
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
// @Description  API for create user
// @Tags         Task
// @Security     Bearer
// @Body 	     request.ArticleRequest
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.UserRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err := _i.userService.Store(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully created"},
	})
}

// Update user
// @Summary      update user
// @Description  API for update user
// @Tags         Task
// @Security     Bearer
// @Body 	     request.ArticleRequest
// @Param        id path int true "User ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
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

	err = _i.userService.Update(id, *req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully updated"},
	})
}

// Delete user
// @Summary      delete user
// @Description  API for delete user
// @Tags         Task
// @Security     Bearer
// @Param        id path int true "User ID"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.userService.Destroy(id)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"User successfully deleted"},
	})
}
