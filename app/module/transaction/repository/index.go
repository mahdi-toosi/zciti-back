package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/transaction/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"

	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.Transactions) (transactions []*schema.Transaction, totalAmount uint64, paging paginator.Pagination, err error)
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

func (_i *repo) GetAll(req request.Transactions) (transactions []*schema.Transaction, totalAmount uint64, paging paginator.Pagination, err error) {
	baseQuery := _i.DB.Main.
		Model(&schema.Transaction{}).
		Where(&schema.Transaction{WalletID: req.WalletID})

	if req.StartTime != nil && !req.StartTime.IsZero() {
		baseQuery = baseQuery.Where("created_at >= ?", req.StartTime)
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		baseQuery = baseQuery.Where("created_at <= ?", req.EndTime)
	}

	if req.Status != "" {
		baseQuery = baseQuery.Where("status = ?", req.Status)
	}

	// Apply filters (CityID / WorkspaceID / DormitoryID)
	if req.CityID > 0 || req.WorkspaceID > 0 || req.DormitoryID > 0 {
		var filters []uint64
		if req.DormitoryID > 0 {
			filters = append(filters, req.DormitoryID)
		} else if req.WorkspaceID > 0 {
			filters = append(filters, req.WorkspaceID)
		} else if req.CityID > 0 {
			filters = append(filters, req.CityID)
		}

		baseQuery = baseQuery.
			Joins("JOIN order_items ON transactions.order_id = order_items.order_id").
			Joins("JOIN posts_taxonomies ON posts_taxonomies.post_id = order_items.post_id").
			Where("posts_taxonomies.taxonomy_id IN (?)", filters)
	}

	// ✅ Clone safely for SUM
	sumQuery := baseQuery.Session(&gorm.Session{}) // new session, same model and conditions
	if err = sumQuery.
		Where(&schema.Transaction{Status: schema.TransactionStatusSuccess}).
		Select("COALESCE(CAST(ROUND(SUM(amount)) AS BIGINT), 0)").
		Scan(&totalAmount).Error; err != nil {
		return
	}

	// ✅ Build list query (add group/pagination only here)
	listQuery := baseQuery.Group("transactions.id")

	if req.Pagination.Page > 0 {
		var total int64
		listQuery.Count(&total)
		req.Pagination.Total = total

		listQuery = listQuery.Offset(req.Pagination.Offset).Limit(req.Pagination.Limit)
	}

	err = listQuery.
		Preload("User").
		Order("created_at desc").
		Find(&transactions).Error
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
