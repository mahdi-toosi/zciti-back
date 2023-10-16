package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/auth/request"
	"go-fiber-starter/app/module/auth/service"
	usersRepo "go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/response"
	"strconv"
)

type IRestController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

func RestController(service service.IService, usersRepo usersRepo.IRepository) IRestController {
	return &controller{service, usersRepo}
}

type controller struct {
	service   service.IService
	usersRepo usersRepo.IRepository
}

// Login
// @Summary      Do log in
// @Description  API for do log in
// @Tags         Authentication
// @Param 		 user body request.Login true "User details"
// @Router       /auth/login [post]
func (_i *controller) Login(c *fiber.Ctx) error {
	req := new(request.Login)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := _i.service.Login(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{Data: res})
}

// Register
// @Summary      register
// @Description  API for register
// @Tags         Authentication
// @Param 		 user body userReq.UserRequest true "User details"
// @Router       /auth/register [post]
func (_i *controller) Register(c *fiber.Ctx) error {
	req := new(request.Register)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	res, err := _i.service.Register(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{Data: res})
}
