package service

import (
	"go-fiber-starter/app/database/schema"
	postService "go-fiber-starter/app/module/post/service"
	"go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/product/request"
	"go-fiber-starter/app/module/product/response"
	uresponse "go-fiber-starter/app/module/user/response"
	userService "go-fiber-starter/app/module/user/service"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.ProductsRequest, isForUser bool) (products []*response.Product, paging paginator.Pagination, err error)
	Show(businessID uint64, id uint64) (product *response.Product, err error)
	Store(req request.Product) (postID *uint64, err error)
	StoreVariant(req request.ProductInPost) (productID *uint64, err error)
	StoreAttribute(req request.StoreProductAttribute) error
	Update(id uint64, req request.Product) (err error)
	Delete(businessID uint64, id uint64) error
	DeleteVariant(businessID uint64, productID uint64, variantID uint64) error
}

func Service(repo repository.IRepository, pService postService.IService, uService userService.IService) IService {
	return &service{
		repo, pService, uService,
	}
}

type service struct {
	Repo     repository.IRepository
	PService postService.IService
	uService userService.IService
}

func (_i *service) Index(req request.ProductsRequest, isForUser bool) (products []*response.Product, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req, isForUser)
	if err != nil {
		return
	}

	var observers []*uresponse.User

	for _, result := range results {
		products = append(products, response.FromDomain(result, result.Products, observers, isForUser))
	}

	return
}

func (_i *service) Show(businessID uint64, id uint64) (product *response.Product, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	observers, err := _i.uService.GetPostObservers(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result, result.Products, observers, false), nil
}

func (_i *service) Store(req request.Product) (postID *uint64, err error) {
	post, err := _i.PService.Store(req.Post)
	if err != nil {
		return nil, err
	}

	products := req.ToDomain(post.ID, req.Post.BusinessID)

	if err = _i.Repo.Creates(products); err != nil {
		return nil, err
	}

	return &post.ID, nil
}

func (_i *service) StoreVariant(req request.ProductInPost) (productID *uint64, err error) {
	product := req.ToDomain(req.PostID, req.BusinessID)

	if product.ID == 0 {
		if err = _i.Repo.Create(product); err != nil {
			return nil, err
		}
	} else {
		if err = _i.Repo.Update(product); err != nil {
			return nil, err
		}
	}
	return &product.ID, nil
}

func (_i *service) StoreAttribute(req request.StoreProductAttribute) (err error) {
	if err = _i.Repo.CreateAttribute(req.ProductID, req.AddedAttrID); err != nil {
		return err
	}

	if req.RemovedAttrID != 0 {
		if err = _i.Repo.DeleteAttribute(req.ProductID, req.RemovedAttrID); err != nil {
			return err
		}
	}

	return nil
}

func (_i *service) Update(id uint64, req request.Product) (err error) {
	if err = _i.PService.Update(id, req.Post); err != nil {
		return err
	}

	products := req.ToDomain(id, req.Post.BusinessID)

	var createList []*schema.Product
	var updateList []*schema.Product
	for _, p := range products {
		if p.ID == 0 {
			createList = append(createList, p)
		} else {
			updateList = append(updateList, p)
		}
	}
	if len(createList) > 0 {
		if err = _i.Repo.Creates(createList); err != nil {
			return err
		}
	}
	return _i.Repo.Updates(updateList)
}

func (_i *service) Delete(businessID uint64, id uint64) error {
	return _i.Repo.Delete(businessID, id)
}

func (_i *service) DeleteVariant(businessID uint64, productID uint64, variantID uint64) error {
	return _i.Repo.DeleteVariant(businessID, productID, variantID)
}
