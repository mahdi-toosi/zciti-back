package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/taxonomy/request"
	"go-fiber-starter/app/module/taxonomy/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"strconv"
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

// Index
// @Summary      Get all taxonomies
// @Tags         Taxonomies
// @Security     Bearer
// @Router       /taxonomies [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Taxonomies
	req.Pagination = paginate
	req.Keyword = c.Query("Keyword")

	taxonomies, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: taxonomies,
		Meta: paging,
	})
}

// Show
// @Summary      Get one taxonomy
// @Tags         Taxonomies
// @Security     Bearer
// @Param        id path int true "Taxonomy ID"
// @Router       /taxonomies/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	taxonomy, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return c.JSON(taxonomy)
}

// Store
// @Summary      Create taxonomy
// @Tags         Taxonomies
// @Param 		 taxonomy body request.Taxonomy true "Taxonomy details"
// @Router       /taxonomies [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.Taxonomy)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update
// @Summary      update taxonomy
// @Security     Bearer
// @Tags         Taxonomies
// @Param 		 taxonomy body request.Taxonomy true "Taxonomy details"
// @Param        id path int true "Taxonomy ID"
// @Router       /taxonomies/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Taxonomy)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete
// @Summary      delete taxonomy
// @Tags         Taxonomies
// @Security     Bearer
// @Param        id path int true "Taxonomy ID"
// @Router       /taxonomies/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
