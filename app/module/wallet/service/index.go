package service

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/wallet/repository"
	"go-fiber-starter/app/module/wallet/request"
	"go-fiber-starter/app/module/wallet/response"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IService interface {
	Index(req request.Wallets) (wallets []*response.Wallet, paging paginator.Pagination, err error)
	Show(id *uint64, userID *uint64, businessID *uint64) (wallet *response.Wallet, err error)
	Store(req *schema.Wallet, tx *gorm.DB) (err error)
	Update(id uint64, req schema.Wallet) (err error)
	Destroy(id uint64) error
	GetOrCreateWallet(userID *uint64, businessID *uint64, tx *gorm.DB) (wallet *response.Wallet, err error)
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.Wallets) (wallets []*response.Wallet, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		wallets = append(wallets, response.FromDomain(result))
	}
	return
}

func (_i *service) Show(id *uint64, userID *uint64, businessID *uint64) (article *response.Wallet, err error) {
	result, err := _i.Repo.GetOne(id, userID, businessID)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (_i *service) Store(req *schema.Wallet, tx *gorm.DB) (err error) {
	return _i.Repo.Create(req, tx)
}

func (_i *service) Update(id uint64, req schema.Wallet) (err error) {
	return _i.Repo.Update(id, &req)
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}

func (_i *service) GetOrCreateWallet(userID *uint64, businessID *uint64, tx *gorm.DB) (wallet *response.Wallet, err error) {
	wallet, err = _i.Show(nil, userID, businessID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var walletSchema = schema.Wallet{}
			if userID != nil {
				walletSchema.UserID = userID
			}
			if businessID != nil {
				walletSchema.BusinessID = businessID
			}
			err := _i.Store(&walletSchema, tx)
			if err != nil {
				return nil, err
			}

			wallet = &response.Wallet{
				ID:     walletSchema.ID,
				Amount: walletSchema.Amount,
			}
		} else {
			return nil, err
		}
	}

	return wallet, nil
}
