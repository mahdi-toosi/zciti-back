package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/product/request"
	"go-fiber-starter/app/module/product/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	StoreVariant(c *fiber.Ctx) error
	StoreAttribute(c *fiber.Ctx) error
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
// @Summary      Get all products
// @Tags         Product
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/products [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.ProductsRequest
	req.Pagination = paginate
	req.BusinessID = businessID
	req.Keyword = c.Query("keyword")
	req.CategoryID = c.Query("CategoryID")

	products, paging, err := _i.service.Index(req, utils.IsForUser(c))
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: products,
		Meta: paging,
	})
}

// Show
// @Summary      Get one product
// @Tags         Product
// @Security     Bearer
// @Param        id path int true "Product ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/products/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	product, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(product)
}

// Store
// @Summary      Create product
// @Tags         Product
// @Param 		 product body request.Product true "Product details"
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/products [product]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrForbidden
	}

	req := new(request.Product)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.Post.AuthorID = user.ID
	req.Post.BusinessID = businessID

	p, err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

// StoreVariant
// @Summary      Create product
// @Tags         Product
// @Param 		 product body request.ProductInPost true "Product details"
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/products/:id/product-variant [product]
func (_i *controller) StoreVariant(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	postID, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.ProductInPost)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.PostID = postID
	req.BusinessID = businessID

	p, err := _i.service.StoreVariant(*req)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

// StoreAttribute
// @Summary      Create product
// @Tags         Product
// @Param 		 product body request.ProductInPost true "Product details"
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Param        id path int true "Product ID"
// @Router       /business/:businessID/products/:id/product-attribute [product]
func (_i *controller) StoreAttribute(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	productID, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.StoreProductAttribute)
	req.ProductID = productID
	req.BusinessID = businessID

	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := _i.service.StoreAttribute(*req); err != nil {
		return err
	}

	return c.JSON("success")
}

// Update
// @Summary      update product
// @Tags         Product
// @Param 		 product body request.Product true "Product details"
// @Security     Bearer
// @Param        id path int true "Product ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/products/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Product)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.Post.BusinessID = businessID
	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete
// @Summary      delete product
// @Tags         Product
// @Security     Bearer
// @Param        id path int true "Product ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/products/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	err = _i.service.Delete(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
