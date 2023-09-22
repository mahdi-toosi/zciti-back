package service

import (
	"go-fiber-starter/app/module/post/repository"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/app/module/post/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.PostsRequest) (posts []*response.Post, paging paginator.Pagination, err error)
	Show(id uint64) (post *response.Post, err error)
	Store(req request.PostRequest) (err error)
	Update(id uint64, req request.PostRequest) (err error)
	Delete(id uint64) error
}

func Service(repo repository.IRepository) IService {
	return &service{
		repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.PostsRequest) (posts []*response.Post, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		posts = append(posts, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(id uint64) (post *response.Post, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.PostRequest) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.PostRequest) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Delete(id uint64) error {
	return _i.Repo.Delete(id)
}
