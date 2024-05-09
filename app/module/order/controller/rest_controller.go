package controller

import (
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/order/service"
	urequest "go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	StoreUniWash(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index all Orders
// @Summary      Get all orders
// @Tags         Orders
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Orders
	req.BusinessID = businessID
	req.Pagination = paginate

	orders, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: orders,
		Meta: paging,
	})
}

// Show one Order
// @Summary      Get one order
// @Tags         Orders
// @Security     Bearer
// @Param        id path int true "Order ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	order, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(order)
}

// Store order
// @Summary      Create order
// @Tags         Orders
// @Param 		 order body request.Order true "Order details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.Order)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.BusinessID = businessID
	_, err = _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// StoreUniWash order
// @Summary      Create order
// @Tags         Orders
// @Param 		 order body request.Order true "Order details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders [post]
func (_i *controller) StoreUniWash(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	req := new(urequest.StoreUniWash)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.UserID = user.ID
	req.BusinessID = businessID
	err = _i.service.StoreUniWash(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update order
// @Summary      update order
// @Security     Bearer
// @Tags         Orders
// @Param 		 order body request.Order true "Order details"
// @Param        id path int true "Order ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Order)
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

// Delete order
// @Summary      delete order
// @Tags         Orders
// @Security     Bearer
// @Param        id path int true "Order ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders/:id [delete]
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
