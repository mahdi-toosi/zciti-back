package repository

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/wallet/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.Wallets) (wallets []*schema.Wallet, paging paginator.Pagination, err error)
	GetOne(id *uint64, userID *uint64, businessID *uint64) (wallet *schema.Wallet, err error)
	Create(wallet *schema.Wallet, tx *gorm.DB) (err error)
	Update(id uint64, wallet *schema.Wallet) (err error)
	Delete(id uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Wallets) (wallets []*schema.Wallet, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Wallet{})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&wallets).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id *uint64, userID *uint64, businessID *uint64) (wallet *schema.Wallet, err error) {
	// Start building the query
	query := _i.DB.Main.Model(&schema.Wallet{})

	// Add conditions based on which argument is provided
	if id != nil {
		query = query.Where("id = ?", id)
	} else if userID != nil {
		query = query.Where("user_id = ?", userID)
	} else if businessID != nil {
		query = query.Where("business_id = ?", businessID)
	} else {
		return nil, errors.New("no id or userID or businessID provided")
	}

	// Execute the query
	if err := query.First(&wallet).Error; err != nil {
		return nil, err
	}

	return wallet, nil
}

func (_i *repo) Create(wallet *schema.Wallet, tx *gorm.DB) (err error) {
	if tx != nil {
		err = tx.Create(&wallet).Error
	} else {
		err = _i.DB.Main.Create(&wallet).Error
	}
	return err
}

func (_i *repo) Update(id uint64, wallet *schema.Wallet) (err error) {
	return _i.DB.Main.Model(&schema.Wallet{}).
		Where(&schema.Wallet{ID: id, BusinessID: wallet.BusinessID}).
		Updates(wallet).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Wallet{}, id).Error
}
