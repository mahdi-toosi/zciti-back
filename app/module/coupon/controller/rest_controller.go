package controller

import (
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/coupon/request"
	"go-fiber-starter/app/module/coupon/service"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error

	ValidateCoupon(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index all Coupons
// @Summary      Get all coupons
// @Tags         Coupons
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/coupons [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Coupons
	req.BusinessID = businessID
	req.Pagination = paginate
	req.Title = c.Query("Title")

	coupons, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: coupons,
		Meta: paging,
	})
}

// Show one Coupon
// @Summary      Get one coupon
// @Tags         Coupons
// @Security     Bearer
// @Param        id path int true "Coupon ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/coupons/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	coupon, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(coupon)
}

// Store coupon
// @Summary      Create coupon
// @Tags         Coupons
// @Param 		 coupon body request.Coupon true "Coupon details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/coupons [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.Coupon)
	req.BusinessID = businessID

	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}
	err = _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update coupon
// @Summary      update coupon
// @Security     Bearer
// @Tags         Coupons
// @Param 		 coupon body request.Coupon true "Coupon details"
// @Param        id path int true "Coupon ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/coupons/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Coupon)
	req.BusinessID = businessID
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete coupon
// @Summary      delete coupon
// @Tags         Coupons
// @Security     Bearer
// @Param        id path int true "Coupon ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/coupons/:id [delete]
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

// ValidateCoupon
// @Summary      Validate coupon
// @Tags         Coupons
// @Param 		 coupon body request.ValidateCoupon true "Coupon details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/validate-coupon [post]
func (_i *controller) ValidateCoupon(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.ValidateCoupon)
	req.BusinessID = businessID

	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	coupon, err := _i.service.ValidateCoupon(*req)
	if err != nil {
		return err
	}

	totalAmt := _i.service.CalcTotalAmtWithDiscount(coupon, &req.OrderTotalAmt)

	return response.Resp(c, response.Response{
		Data: totalAmt,
	})
}
