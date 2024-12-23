package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-fiber-starter/app/module/asset/request"
	assetsResponse "go-fiber-starter/app/module/asset/response"
	"go-fiber-starter/app/module/asset/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"path/filepath"
	"strings"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index
// @Summary      Get all assets
// @Tags         Assets
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/assets [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Assets
	req.Pagination = paginate
	req.BusinessID = businessID
	req.Keyword = c.Query("Keyword")

	assets, assetsSize, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: assets,
		Meta: assetsResponse.Assets{
			Meta:       paging,
			AssetsSize: assetsSize,
		},
	})
}

// Store
// @Summary      Create asset
// @Tags         Assets
// @Param 		 asset body request.Asset true "Asset details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/assets [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	asset, err := c.FormFile("Asset")
	if err != nil {
		return err
	}

	req := new(request.Asset)
	req.Asset = *asset
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.Ext = strings.TrimPrefix(filepath.Ext(req.Asset.Filename), ".")

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	req.UserID = user.ID
	req.Title = strings.TrimSuffix(req.Asset.Filename, filepath.Ext(req.Asset.Filename))

	err = _i.service.Store(c, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete
// @Summary      delete asset
// @Tags         Assets
// @Security     Bearer
// @Param        id path int true "Asset ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/assets/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	err = _i.service.Destroy(user, id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
