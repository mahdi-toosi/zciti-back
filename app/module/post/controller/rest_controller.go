package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/app/module/post/service"
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
	DeleteTaxonomies(c *fiber.Ctx) error
	InsertTaxonomies(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index
// @Summary      Get all posts
// @Tags         Post
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       business/:businessID/posts [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.PostsRequest
	req.Pagination = paginate
	req.BusinessID = businessID
	req.Keyword = c.Query("keyword")

	posts, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: posts,
		Meta: paging,
	})
}

// Show
// @Summary      Get one post
// @Tags         Post
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/posts/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	post, err := _i.service.Show(businessID, id)
	if err != nil {
		return err
	}

	return c.JSON(post)
}

// Store
// @Summary      Create post
// @Tags         Post
// @Param 		 post body request.Post true "Post details"
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/posts [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrForbidden
	}

	req := new(request.Post)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.AuthorID = user.ID
	req.BusinessID = businessID

	p, err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

// Update
// @Summary      update post
// @Tags         Post
// @Param 		 post body request.Post true "Post details"
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/posts/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Post)
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
// @Summary      delete post
// @Tags         Post
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/posts/:id [delete]
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

// DeleteTaxonomies
// @Summary      delete post
// @Tags         Post
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/posts/:id/delete-taxonomies [delete]
func (_i *controller) DeleteTaxonomies(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.PostTaxonomies)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.PostID = id
	req.BusinessID = businessID
	err = _i.service.DeleteTaxonomies(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// InsertTaxonomies
// @Summary      delete post
// @Tags         Post
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/posts/:id/insert-taxonomies [delete]
func (_i *controller) InsertTaxonomies(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.PostTaxonomies)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.PostID = id
	req.BusinessID = businessID
	err = _i.service.InsertTaxonomies(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}
