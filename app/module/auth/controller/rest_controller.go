package controller

import (
	"github.com/bangadam/go-fiber-starter/app/module/auth/request"
	"github.com/bangadam/go-fiber-starter/app/module/auth/service"
	"github.com/bangadam/go-fiber-starter/utils/response"
	"github.com/gofiber/fiber/v2"
)

type IRestController interface {
	Login(c *fiber.Ctx) error
}

func RestController(service service.IService) IRestController {
	return &controller{service}
}

type controller struct {
	service service.IService
}

// Login
// @Summary      Do log in
// @Description  API for do log in
// @Tags         Authentication
// @Security     Bearer
// @Body 	     request.LoginRequest
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      422  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/login [post]
func (_i *controller) Login(c *fiber.Ctx) error {
	req := new(request.LoginRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := _i.service.Login(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Login success"},
		Code:     fiber.StatusOK,
	})
}
