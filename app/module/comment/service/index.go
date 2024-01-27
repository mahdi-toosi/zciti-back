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
	Show(id uint64) (comment *response.Comment, err error)
	Store(postID uint64, req request.Comment) (err error)
	Update(postID uint64, id uint64, req request.Comment) (err error)
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

func (_i *service) Show(id uint64) (comment *response.Comment, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(postID uint64, req request.Comment) (err error) {
	log.Debug().Msgf(" req.AuthorID => %+v ", req.AuthorID)
	if ok, _ := _i.hasPermissionToStore(postID, req.AuthorID); !ok {
		return fiber.ErrForbidden
	}

	return _i.Repo.Create(req.ToDomain(&postID))
}

func (_i *service) Update(postID uint64, id uint64, req request.Comment) (err error) {
	if ok, _ := _i.hasPermissionToStore(postID, req.AuthorID); !ok {
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

func (_i *service) hasPermissionToStore(postID uint64, authorID uint64) (ok bool, err error) {
	post, err := _i.PostsRepo.GetOne(postID)
	if err != nil {
		return false, nil
	}

	if post.CommentsStatus == schema.PostCommentStatusClose {
		return false, nil
	}

	if post.CommentsStatus == schema.PostCommentStatusOpen {
		return true, nil
	}

	if post.CommentsStatus == schema.PostCommentStatusOnlyCustomers {
		log.Debug().Msgf("%+v", authorID)
		// TODO => complete it
	}

	if post.CommentsStatus == schema.PostCommentStatusOnlyBuyers {
		log.Debug().Msgf("%+v", authorID)
		// TODO => complete it
	}

	return true, nil
}
