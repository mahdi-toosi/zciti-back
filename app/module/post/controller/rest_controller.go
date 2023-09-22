package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/app/module/post/service"
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

// Index all Posts
// @Summary      Get all posts
// @Tags         Task
// @Security     Bearer
// @Router       /posts [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.PostsRequest
	req.Pagination = paginate

	posts, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Post list successfully retrieved"},
		Data:     posts,
		Meta:     paging,
	})
}

// Show one Post
// @Summary      Get one post
// @Tags         Task
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Router       /posts/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	posts, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Post successfully retrieved"},
		Data:     posts,
	})
}

// Store post
// @Summary      Create post
// @Tags         Task
// @Param 		 post body request.PostRequest true "Post details"
// @Security     Bearer
// @Router       /posts [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.PostRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Post successfully created"},
	})
}

// Update post
// @Summary      update post
// @Tags         Task
// @Param 		 post body request.PostRequest true "Post details"
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Router       /posts/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	req := new(request.PostRequest)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Post successfully updated"},
	})
}

// Delete post
// @Summary      delete post
// @Tags         Task
// @Security     Bearer
// @Param        id path int true "Post ID"
// @Router       /posts/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.service.Delete(id)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Post successfully deleted"},
	})
}
