package controller

import (
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/example/request"
	"go-fiber-starter/app/module/example/service"
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

// Index all Examples
// @Summary      Get all examples
// @Tags         Examples
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/examples [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Examples
	req.BusinessID = businessID
	req.Pagination = paginate

	examples, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: examples,
		Meta: paging,
	})
}

// Show one Example
// @Summary      Get one example
// @Tags         Examples
// @Security     Bearer
// @Param        id path int true "Example ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/examples/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	example, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(example)
}

// Store example
// @Summary      Create example
// @Tags         Examples
// @Param 		 example body request.Example true "Example details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/examples [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.Example)
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

// Update example
// @Summary      update example
// @Security     Bearer
// @Tags         Examples
// @Param 		 example body request.Example true "Example details"
// @Param        id path int true "Example ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/examples/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Example)
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

// Delete example
// @Summary      delete example
// @Tags         Examples
// @Security     Bearer
// @Param        id path int true "Example ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/examples/:id [delete]
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
