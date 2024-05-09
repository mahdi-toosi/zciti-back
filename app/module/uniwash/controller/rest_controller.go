package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/app/module/uniwash/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	SendCommand(c *fiber.Ctx) error
	IndexReservedMachines(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// SendCommand
// @Summary      Send command to washing machines
// @Tags         UniWash
// @Param 		 taxonomy body request.SendCommand true "SendCommand details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/uni-wash/send-command [post]
// @Router       /user/business/:businessID/uni-wash/send-command [post]
func (_i *controller) SendCommand(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	req := new(request.SendCommand)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.UserID = user.ID
	req.BusinessID = businessID
	err = _i.service.SendCommand(*req, utils.IsForUser(c))
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// IndexReservedMachines
// @Summary      Reservations list
// @Tags         UniWash
// @Param 		 taxonomy body request.ReservedMachinesRequest true "Reservation params"
// @Param        businessID path int true "Business ID"
// @Router       /user/business/:businessID/uni-wash/reserved-machines [post]
func (_i *controller) IndexReservedMachines(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	req := new(request.ReservedMachinesRequest)
	req.Pagination = paginate
	req.BusinessID = businessID
	if utils.IsForUser(c) {
		req.UserID = user.ID
	}

	reservations, paging, err := _i.service.IndexReservedMachines(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: reservations,
		Meta: paging,
	})
}
