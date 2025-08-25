package controller

import (
	"go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/app/module/uniwash/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IRestController interface {
	SendCommand(c *fiber.Ctx) error
	IndexReservedMachines(c *fiber.Ctx) error
	CheckLastCommandStatus(c *fiber.Ctx) error
	GetReservationOptions(c *fiber.Ctx) error
	SendDeviceIsOffMsgToUser(c *fiber.Ctx) error
	SendFullCouponToUser(c *fiber.Ctx) error
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

	if c.Query("Date") != "" {
		loc, _ := time.LoadLocation("Asia/Tehran")
		date, err := time.ParseInLocation(time.DateOnly, c.Query("Date"), loc)
		if err != nil {
			return err
		}
		req.Date = date
	}

	req.With = c.Query("With")
	if utils.IsForUser(c) {
		req.UserID = user.ID
	}
	if c.Query("ProductID") != "" {
		pID, err := utils.GetUintInQueries(c, "ProductID")
		if err != nil {
			return err
		}
		req.ProductID = pID
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

// CheckLastCommandStatus
// @Summary      Reservations list
// @Tags         UniWash
// @Param        businessID path int true "Business ID"
// @Param        reservationID path int true "Reservation ID"
// @Router       /business/:businessID/uni-wash/check-last-command-status/:reservationID [get]
func (_i *controller) CheckLastCommandStatus(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	reservationID, err := utils.GetIntInParams(c, "reservationID")
	if err != nil {
		return err
	}

	status, err := _i.service.CheckLastCommandStatus(businessID, reservationID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: status,
	})
}

// GetReservationOptions
// @Summary      Device reservations options
// @Tags         UniWash
// @Param        businessID path int true "Business ID"
// @Param        reservationID path int true "Reservation ID"
// @Router       /business/:businessID/uni-wash/device/reservation-options [get]
func (_i *controller) GetReservationOptions(c *fiber.Ctx) error {
	return response.Resp(c, response.Response{
		Data: _i.service.GetReservationOptions(),
	})
}

// SendDeviceIsOffMsgToUser
// @Summary      Send device is off msg to user
// @Tags         UniWash
// @Param        businessID path int true "Business ID"
// @Param        reservationID path int true "Reservation ID"
// @Router       /business/:businessID/uni-wash/send-device-is-off-msg-to-user/:reservationID [get]
func (_i *controller) SendDeviceIsOffMsgToUser(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	reservationID, err := utils.GetIntInParams(c, "reservationID")
	if err != nil {
		return err
	}

	err = _i.service.SendDeviceIsOffMsgToUser(businessID, reservationID)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// SendFullCouponToUser
// @Summary      Send device is off msg to user
// @Tags         UniWash
// @Param        businessID path int true "Business ID"
// @Param        reservationID path int true "Reservation ID"
// @Router       /business/:businessID/uni-wash/send-full-coupon-to-user/:reservationID [get]
func (_i *controller) SendFullCouponToUser(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	reservationID, err := utils.GetIntInParams(c, "reservationID")
	if err != nil {
		return err
	}

	err = _i.service.SendFullCouponToUser(businessID, reservationID)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
