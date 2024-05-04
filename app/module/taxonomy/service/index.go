package service

import (
	"go-fiber-starter/app/module/taxonomy/repository"
	"go-fiber-starter/app/module/taxonomy/request"
	"go-fiber-starter/app/module/taxonomy/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Taxonomies, forUser bool) (taxonomies []*response.Taxonomy, paging paginator.Pagination, err error)
	Search(req request.Taxonomies, forUser bool) (taxonomies []*response.Taxonomy, paging paginator.Pagination, err error)
	Show(BusinessID uint64, id uint64) (taxonomy *response.Taxonomy, err error)
	Store(req request.Taxonomy) (err error)
	Update(id uint64, req request.Taxonomy) (err error)
	Destroy(BusinessID uint64, id uint64) error
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.Taxonomies, forUser bool) (taxonomies []*response.Taxonomy, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		taxonomies = append(taxonomies, response.FromDomain(result, forUser))
	}

	return
}

func (_i *service) Search(req request.Taxonomies, forUser bool) (taxonomies []*response.Taxonomy, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.Search(req)
	if err != nil {
		return
	}

	for _, result := range results {
		taxonomies = append(taxonomies, response.FromDomain(result, forUser))
	}

	return
}
func (_i *service) Show(businessID uint64, id uint64) (taxonomy *response.Taxonomy, err error) {
	result, err := _i.Repo.GetOne(businessID, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result, false), nil
}

func (_i *service) Store(req request.Taxonomy) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.Taxonomy) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(businessID uint64, id uint64) error {
	return _i.Repo.Delete(businessID, &id)
}
