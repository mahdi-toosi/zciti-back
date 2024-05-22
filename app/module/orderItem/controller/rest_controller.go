package controller

import (
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/app/module/orderItem/service"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	//Store(c *fiber.Ctx) error
	//Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index all OrderItems
// @Summary      Get all orderItems
// @Tags         OrderItems
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orderItems [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.OrderItems
	req.Pagination = paginate

	orderItems, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: orderItems,
		Meta: paging,
	})
}

// Show one OrderItem
// @Summary      Get one orderItem
// @Tags         OrderItems
// @Security     Bearer
// @Param        id path int true "OrderItem ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orderItems/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	orderItem, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return c.JSON(orderItem)
}

//// Store orderItem
//// @Summary      Create orderItem
//// @Tags         OrderItems
//// @Param 		 orderItem body request.OrderItem true "OrderItem details"
//// @Param        businessID path int true "Business ID"
//// @Router       /business/:businessID/orderItems [post]
//func (_i *controller) Store(c *fiber.Ctx) error {
//	req := new(request.OrderItem)
//	if err := response.ParseAndValidate(c, req); err != nil {
//		return err
//	}
//
//	err := _i.service.Store(*req)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON("success")
//}

//// Update orderItem
//// @Summary      update orderItem
//// @Security     Bearer
//// @Tags         OrderItems
//// @Param 		 orderItem body request.OrderItem true "OrderItem details"
//// @Param        id path int true "OrderItem ID"
//// @Param        businessID path int true "Business ID"
//// @Router       /business/:businessID/orderItems/:id [put]
//func (_i *controller) Update(c *fiber.Ctx) error {
//	id, err := utils.GetIntInParams(c, "id")
//	if err != nil {
//		return err
//	}
//
//	req := new(request.OrderItem)
//	if err := response.ParseAndValidate(c, req); err != nil {
//		return err
//	}
//
//	err = _i.service.Update(id, *req)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON("success")
//}

// Delete orderItem
// @Summary      delete orderItem
// @Tags         OrderItems
// @Security     Bearer
// @Param        id path int true "OrderItem ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orderItems/:id [delete]
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
