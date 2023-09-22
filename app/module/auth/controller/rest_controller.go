package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/auth/request"
	"go-fiber-starter/app/module/auth/service"
	"go-fiber-starter/utils/response"
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
// @Param 		 user body request.LoginRequest true "User details"
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
