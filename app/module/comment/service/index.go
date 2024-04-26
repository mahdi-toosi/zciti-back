package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/comment/repository"
	"go-fiber-starter/app/module/comment/request"
	"go-fiber-starter/app/module/comment/response"
	postsRepository "go-fiber-starter/app/module/post/repository"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(postID uint64, req request.Comments) (comments []*response.Comment, paging paginator.Pagination, err error)
	Show(id uint64) (comment *schema.Comment, err error)
	Store(businessID uint64, postID uint64, req request.Comment, isUserOwnerOfSomeBusiness bool) (err error)
	Update(businessID uint64, postID uint64, id uint64, req request.Comment) (err error)
	UpdateStatus(id uint64, req request.UpdateCommentStatus) (err error)
	Destroy(id uint64) error
}

func Service(Repo repository.IRepository, PostsRepo postsRepository.IRepository) IService {
	return &service{
		Repo,
		PostsRepo,
	}
}

type service struct {
	Repo      repository.IRepository
	PostsRepo postsRepository.IRepository
}

func (_i *service) Index(postID uint64, req request.Comments) (comments []*response.Comment, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(postID, req)
	if err != nil {
		return
	}

	for _, result := range results {
		comments = append(comments, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(id uint64) (comment *schema.Comment, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (_i *service) Store(businessID uint64, postID uint64, req request.Comment, isUserOwnerOfSomeBusiness bool) (err error) {
	ok, _, post := _i.hasPermissionToStore(businessID, postID, req.AuthorID)
	if !ok {
		return fiber.ErrForbidden
	}

	var IsNotBusinessOwner = false
	req.IsBusinessOwner = &IsNotBusinessOwner

	if isUserOwnerOfSomeBusiness {
		if err == nil && post.Business.OwnerID == req.AuthorID {
			req.Status = schema.CommentStatusApproved
			var IsBusinessOwner = true
			req.IsBusinessOwner = &IsBusinessOwner
		}
	}

	err = _i.Repo.Create(req.ToDomain(&postID))
	if err != nil {
		return err
	}

	err = _i.PostsRepo.UpdateCommentCount(postID, "+ 1")
	if err != nil {
		return err
	}

	return
}

func (_i *service) Update(businessID uint64, postID uint64, id uint64, req request.Comment) (err error) {
	if ok, _, _ := _i.hasPermissionToStore(businessID, postID, req.AuthorID); !ok {
		return fiber.ErrForbidden
	}

	return _i.Repo.Update(id, req.ToDomain(nil))
}

func (_i *service) UpdateStatus(id uint64, req request.UpdateCommentStatus) (err error) {
	return _i.Repo.Update(id, &schema.Comment{Status: req.Status})
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}

func (_i *service) hasPermissionToStore(businessID uint64, postID uint64, authorID uint64) (ok bool, err error, post *schema.Post) {
	post, err = _i.PostsRepo.GetOne(businessID, postID)
	if err != nil {
		return false, nil, nil
	}

	if post.Meta.CommentsStatus == schema.PostCommentStatusClose {
		return false, nil, post
	}

	if post.Meta.CommentsStatus == schema.PostCommentStatusOpen {
		return true, nil, post
	}

	if post.Meta.CommentsStatus == schema.PostCommentStatusOnlyBusinessUsers {
		log.Debug().Msgf("%+v", authorID)
		// TODO => complete it
	}

	if post.Meta.CommentsStatus == schema.PostCommentStatusOnlyBuyers {
		log.Debug().Msgf("%+v", authorID)
		// TODO => complete it
	}

	return true, nil, post
}
