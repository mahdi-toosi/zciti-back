package service

import (
	postService "go-fiber-starter/app/module/post/service"
	"go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/product/request"
	"go-fiber-starter/app/module/product/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.ProductsRequest) (products []*response.Product, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (product *response.Product, err error)
	Store(req request.Product) (product *response.Product, err error)
	Update(id uint64, req request.Product) (err error)
	Delete(businessID uint64, id uint64) error
}

func Service(repo repository.IRepository, pService postService.IService) IService {
	return &service{
		repo, pService,
	}
}

type service struct {
	Repo     repository.IRepository
	PService postService.IService
}

func (_i *service) Index(req request.ProductsRequest) (products []*response.Product, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		products = append(products, response.FromDomain(result, result.Products))
	}

	return
}

func (_i *service) Show(businessID uint64, id uint64) (product *response.Product, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result, result.Products), nil
}

func (_i *service) Store(req request.Product) (product *response.Product, err error) {
	post, err := _i.PService.Store(req.Post)
	if err != nil {
		return nil, err
	}

	products := req.ToDomain(post.ID, req.Post.BusinessID)

	err = _i.Repo.Create(products)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (_i *service) Update(id uint64, req request.Product) (err error) {
	err = _i.PService.Update(id, req.Post)
	if err != nil {
		return err
	}

	products := req.ToDomain(id, req.Post.BusinessID)

	return _i.Repo.Update(products)
}

func (_i *service) Delete(businessID uint64, id uint64) error {
	return _i.Repo.Delete(businessID, id)
}
