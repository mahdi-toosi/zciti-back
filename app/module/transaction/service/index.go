package service

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/transaction/repository"
	"go-fiber-starter/app/module/transaction/request"
	"go-fiber-starter/app/module/transaction/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Transactions) (transactions []*response.Transaction, paging paginator.Pagination, err error)
	Show(id uint64) (transaction *response.Transaction, err error)
	Store(req *schema.Transaction) (err error)
	Update(id uint64, req *schema.Transaction) (err error)
	Destroy(id uint64) error
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.Transactions) (transactions []*response.Transaction, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		transactions = append(transactions, response.FromDomain(result))
	}

	return
}

func (_i *service) Show(id uint64) (article *response.Transaction, err error) {
	result, err := _i.Repo.GetOne(&id, nil)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req *schema.Transaction) (err error) {
	return _i.Repo.Create(req, nil)
}

func (_i *service) Update(id uint64, req *schema.Transaction) (err error) {
	// TODO : check business id permission
	return _i.Repo.Update(id, req)
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
