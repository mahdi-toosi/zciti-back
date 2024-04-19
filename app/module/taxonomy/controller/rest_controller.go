package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/taxonomy/request"
	"go-fiber-starter/app/module/taxonomy/service"
	"go-fiber-starter/utils"
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

// Index
// @Summary      Get all taxonomies
// @Tags         Taxonomies
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/taxonomies [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	req := new(request.Taxonomies)
	req.Pagination = paginate
	req.BusinessID = businessID
	req.Keyword = c.Query("Keyword")
	if c.Query("Type") != "" {
		req.Type = schema.TaxonomyType(c.Query("Type"))
	}
	if c.Query("Domain") != "" {
		req.Domain = schema.PostType(c.Query("Domain"))
	}

	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	taxonomies, paging, err := _i.service.Index(*req)
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
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/taxonomies/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	taxonomy, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(taxonomy)
}

// Store
// @Summary      Create taxonomy
// @Tags         Taxonomies
// @Param 		 taxonomy body request.Taxonomy true "Taxonomy details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/taxonomies [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	req := new(request.Taxonomy)
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

// Update
// @Summary      update taxonomy
// @Security     Bearer
// @Tags         Taxonomies
// @Param 		 taxonomy body request.Taxonomy true "Taxonomy details"
// @Param        id path int true "Taxonomy ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/taxonomies/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Taxonomy)
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

// Delete
// @Summary      delete taxonomy
// @Tags         Taxonomies
// @Security     Bearer
// @Param        id path int true "Taxonomy ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/taxonomies/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	err = _i.service.Destroy(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
