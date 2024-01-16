package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/auth/request"
	"go-fiber-starter/app/module/auth/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/response"
	"strconv"
)

type IRestController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	SendOtp(c *fiber.Ctx) error
	ResetPass(c *fiber.Ctx) error
}

func RestController(service service.IService, config *config.Config) IRestController {
	return &controller{service: service, config: config}
}

type controller struct {
	config  *config.Config
	service service.IService
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

	res, err := _i.service.Login(*req, _i.config.Middleware.Jwt)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

// Register
// @Summary      register
// @Description  API for register
// @Tags         Authentication
// @Param 		 user body request.Register true "User details"
// @Router       /auth/register [post]
func (_i *controller) Register(c *fiber.Ctx) error {
	req := new(request.Register)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	res, err := _i.service.Register(req, _i.config.Middleware.Jwt)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

// SendOtp
// @Summary      register
// @Description  API for register
// @Tags         Authentication
// @Param 		 user body request.SendOtp true "User details"
// @Router       /auth/send-otp [post]
func (_i *controller) SendOtp(c *fiber.Ctx) error {
	req := new(request.SendOtp)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	err := _i.service.SendOtp(req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// ResetPass
// @Summary      register
// @Description  API for register
// @Tags         Authentication
// @Param 		 user body request.ResetPass true "User details"
// @Router       /auth/reset-pass [post]
func (_i *controller) ResetPass(c *fiber.Ctx) error {
	req := new(request.ResetPass)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	err := _i.service.ResetPass(req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
