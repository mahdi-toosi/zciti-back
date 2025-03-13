package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/transaction/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.Transactions) (transactions []*schema.Transaction, paging paginator.Pagination, err error)
	GetOne(id *uint64, orderID *uint64) (transaction *schema.Transaction, err error)
	Create(transaction *schema.Transaction, tx *gorm.DB) (err error)
	Update(id uint64, transaction *schema.Transaction) (err error)
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

func (_i *repo) GetAll(req request.Transactions) (transactions []*schema.Transaction, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Transaction{}).
		Where(&schema.Transaction{WalletID: req.WalletID})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("User").Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id *uint64, orderID *uint64) (transaction *schema.Transaction, err error) {
	query := _i.DB.Main.Model(&schema.Transaction{})

	if id != nil {
		query = query.Where("id = ?", id)
	} else if orderID != nil {
		query = query.Where("order_id = ?", orderID)
	}

	if err = query.First(&transaction).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (_i *repo) Create(transaction *schema.Transaction, tx *gorm.DB) (err error) {
	if tx != nil {
		err = tx.Create(&transaction).Error
	} else {
		err = _i.DB.Main.Create(&transaction).Error
	}
	return err
}

func (_i *repo) Update(id uint64, transaction *schema.Transaction) (err error) {
	return _i.DB.Main.Model(&schema.Transaction{}).
		Where(&schema.Transaction{ID: id, WalletID: transaction.WalletID}).
		Updates(transaction).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Transaction{}, id).Error
}
