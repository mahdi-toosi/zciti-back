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
		baseQuery = baseQuery.Where("transactions.created_at >= ?", req.StartTime)
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		baseQuery = baseQuery.Where("transactions.created_at <= ?", req.EndTime)
	}

	if req.Status != "" {
		baseQuery = baseQuery.Where("status = ?", req.Status)
	}

	// Apply filters (CityID / WorkspaceID / DormitoryID)
	if len(req.Taxonomies) > 0 {
		baseQuery = baseQuery.
			Joins("JOIN order_items ON transactions.order_id = order_items.order_id").
			Joins("JOIN posts_taxonomies ON posts_taxonomies.post_id = order_items.post_id").
			Where("posts_taxonomies.taxonomy_id IN (?)", req.Taxonomies)
	}

	// ✅ SUM query with EXISTS to avoid duplicate amounts from JOINs
	sumQuery := _i.DB.Main.
		Model(&schema.Transaction{}).
		Where(&schema.Transaction{WalletID: req.WalletID, Status: schema.TransactionStatusSuccess})

	if req.StartTime != nil && !req.StartTime.IsZero() {
		sumQuery = sumQuery.Where("created_at >= ?", req.StartTime)
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		sumQuery = sumQuery.Where("created_at <= ?", req.EndTime)
	}

	// Use EXISTS subquery to filter without creating duplicate rows
	if len(req.Taxonomies) > 0 {
		sumQuery = sumQuery.Where(
			`EXISTS (
				SELECT 1 FROM order_items 
				JOIN posts_taxonomies ON posts_taxonomies.post_id = order_items.post_id 
				WHERE order_items.order_id = transactions.order_id 
				AND posts_taxonomies.taxonomy_id IN (?)
			)`,
			req.Taxonomies,
		)
	}

	if err = sumQuery.
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
		Order("transactions.created_at desc").
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
