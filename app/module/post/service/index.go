package service

import (
	"go-fiber-starter/app/module/post/repository"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/app/module/post/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.PostsRequest) (posts []*response.Post, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (post *response.Post, err error)
	Store(req request.Post) (post *response.Post, err error)
	Update(id uint64, req request.Post) (err error)
	Delete(businessID uint64, id uint64) error
	DeleteTaxonomies(req request.PostTaxonomies) error
	InsertTaxonomies(req request.PostTaxonomies) error
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

func (_i *service) Show(businessID uint64, id uint64) (post *response.Post, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req request.Post) (post *response.Post, err error) {
	p, err := _i.Repo.Create(req.ToDomain())
	if err != nil {
		return nil, err
	}
	return response.FromDomain(p), nil
}

func (_i *service) Update(id uint64, req request.Post) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Delete(businessID uint64, id uint64) error {
	return _i.Repo.Delete(businessID, id)
}

func (_i *service) DeleteTaxonomies(req request.PostTaxonomies) error {
	return _i.Repo.DeleteTaxonomies(req)
}

func (_i *service) InsertTaxonomies(req request.PostTaxonomies) error {
	return _i.Repo.InsertTaxonomies(req)
}
