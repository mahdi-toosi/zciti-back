package controller

import (
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/app/module/reservation/service"
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

// Index all Reservations
// @Summary      Get all reservations
// @Tags         Reservations
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/reservations [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	var req request.Reservations
	req.Pagination = paginate
	req.BusinessID = businessID
	req.Mobile = c.Query("Mobile")
	req.FullName = c.Query("FullName")
	req.StartTime = utils.GetDateInQueries(c, "StartTime")
	req.ProductID, _ = utils.GetIntInQueries(c, "ProductID")

	if c.Query("UserID") != "" {
		userID, err := utils.GetIntInQueries(c, "UserID")
		if err != nil {
			return err
		}

		if user.ID == userID {
			req.UserID = userID
		}
	}

	reservations, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: reservations,
		Meta: paging,
	})
}

// Show one Reservation
// @Summary      Get one reservation
// @Tags         Reservations
// @Security     Bearer
// @Param        id path int true "Reservation ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/reservations/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	reservation, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(reservation)
}

// Store reservation
// @Summary      Create reservation
// @Tags         Reservations
// @Param 		 reservation body request.Reservation true "Reservation details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/reservations [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.Reservation)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.BusinessID = businessID
	err = _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update reservation
// @Summary      update reservation
// @Security     Bearer
// @Tags         Reservations
// @Param 		 reservation body request.Reservation true "Reservation details"
// @Param        id path int true "Reservation ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/reservations/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Reservation)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.BusinessID = businessID
	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete reservation
// @Summary      delete reservation
// @Tags         Reservations
// @Security     Bearer
// @Param        id path int true "Reservation ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/reservations/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
