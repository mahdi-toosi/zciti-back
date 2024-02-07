package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/comment/request"
	"go-fiber-starter/app/module/comment/service"
	"go-fiber-starter/app/module/post/repository"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	//Delete(c *fiber.Ctx) error
	// Show(c *fiber.Ctx) error
}

func RestController(s service.IService, postsRepo repository.IRepository) IRestController {
	return &controller{s, postsRepo}
}

type controller struct {
	service   service.IService
	postsRepo repository.IRepository
}

// Index
// @Summary      Get all comments
// @Tags         Comments
// @Security     Bearer
// @Router       /posts/:postID/comments [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	postID, err := utils.GetIntInParams(c, "postID")
	if err != nil {
		return err
	}

	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Comments
	req.Pagination = paginate

	comments, paging, err := _i.service.Index(postID, req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: comments,
		Meta: paging,
	})
}

//// Show
//// @Summary      Get one comment
//// @Tags         Comments
//// @Security     Bearer
//// @Param        id path int true "Comment ID"
//// @Router       /posts/:postID/comments/:id [get]
// func (_i *controller) Show(c *fiber.Ctx) error {
//	id, err := utils.GetIntInParams(c, "id")
//	if err != nil {
//		return err
//	}
//
//	comment, err := _i.service.Show(id)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(comment)
//}

// Store
// @Summary      Create comment
// @Tags         Comments
// @Param 		 comment body request.Comment true "Comment details"
// @Router       /posts/:postID/comments [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	postID, err := utils.GetIntInParams(c, "postID")
	if err != nil {
		return err
	}

	req := new(request.Comment)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrBadRequest
	}
	req.AuthorID = user.ID
	req.Status = schema.CommentStatusPending

	err = _i.service.Store(postID, *req, user.IsBusinessOwner())
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update
// @Summary      update comment
// @Security     Bearer
// @Tags         Comments
// @Param 		 comment body request.Comment true "Comment details"
// @Param        id path int true "Comment ID"
// @Router       /posts/:postID/comments/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	postID, err := utils.GetIntInParams(c, "postID")
	if err != nil {
		return err
	}

	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Comment)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrBadRequest
	}

	comment, err := _i.service.Show(id)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if comment.AuthorID != user.ID {
		return fiber.ErrForbidden
	}
	req.Status = schema.CommentStatusPending

	err = _i.service.Update(postID, id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// UpdateStatus
// @Summary      update comment status
// @Security     Bearer
// @Tags         Comments
// @Param 		 comment body request.UpdateCommentStatus true "Comment details"
// @Param        id path int true "Comment ID"
// @Router       /posts/:postID/comments/:id/status [put]
func (_i *controller) UpdateStatus(c *fiber.Ctx) error {
	postID, err := utils.GetIntInParams(c, "postID")
	if err != nil {
		return err
	}

	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.UpdateCommentStatus)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if !user.IsBusinessOwner() {
		return fiber.ErrForbidden
	}
	post, err := _i.postsRepo.GetOne(postID)
	if err != nil || post.Business.OwnerID != user.ID {
		return fiber.ErrBadRequest
	}

	err = _i.service.UpdateStatus(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

//// Delete
//// @Summary      delete comment
//// @Tags         Comments
//// @Security     Bearer
//// @Param        id path int true "Comment ID"
//// @Router       /posts/:postID/comments/:id [delete]
// func (_i *controller) Delete(c *fiber.Ctx) error {
//	id, err := utils.GetIntInParams(c, "id")
//	if err != nil {
//		return err
//	}
//
//	err = _i.service.Destroy(id)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON("success")
//}
