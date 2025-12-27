package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"strings"

	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.Orders) (orders []*schema.Order, totalAmount uint64, paging paginator.Pagination, err error)
	GetOne(userID uint64, id uint64) (order *schema.Order, err error)
	Create(order *schema.Order, tx *gorm.DB) (orderID uint64, err error)
	Update(id uint64, order *schema.Order) (err error)
	Delete(id uint64) (err error)
	BeginTransaction() (*gorm.DB, error)
}

func Repository(db *database.Database) IRepository {
	return &repo{db}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Orders) (orders []*schema.Order, totalAmount uint64, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Order{})

	if req.BusinessID > 0 {
		query = query.Where(&schema.Order{BusinessID: req.BusinessID})
	}

	if req.UserID > 0 {
		query = query.Where(&schema.Order{UserID: req.UserID})
	}

	if req.CouponID > 0 {
		query = query.Where(&schema.Order{CouponID: &req.CouponID})
	}

	// Apply HasCoupon filter
	if req.HasCoupon != nil {
		if *req.HasCoupon {
			query = query.Where("orders.coupon_id IS NOT NULL")
		} else {
			query = query.Where("orders.coupon_id IS NULL")
		}
	}

	if req.Status != "" {
		query = query.Where("orders.status = ?", req.Status)
	}

	// Apply order creation time filters
	if req.OrderStartTime != nil && !req.OrderStartTime.IsZero() {
		query = query.Where("orders.created_at >= ?", utils.StartOfDayString(*req.OrderStartTime))
	}

	if req.OrderEndTime != nil && !req.OrderEndTime.IsZero() {
		query = query.Where("orders.created_at <= ?", utils.EndOfDayString(*req.OrderEndTime))
	}

	// Track if we need joins and grouping
	needsOrderItemsJoin := false
	needsReservationsJoin := false
	needsTaxonomiesJoin := false
	needsUsersJoin := false

	// Check if FullName filter is needed
	if req.FullName != "" {
		needsUsersJoin = true
	}

	// Check if time filters are needed (filter by reservation start/end time)
	if (req.StartTime != nil && !req.StartTime.IsZero()) || (req.EndTime != nil && !req.EndTime.IsZero()) {
		needsOrderItemsJoin = true
		needsReservationsJoin = true
	}

	// Check if ProductID filter is needed
	if req.ProductID != 0 {
		needsOrderItemsJoin = true
	}

	// Check if Taxonomies filter is needed
	if len(req.Taxonomies) > 0 {
		needsOrderItemsJoin = true
		needsTaxonomiesJoin = true
	}

	// Apply joins
	if needsUsersJoin {
		query = query.Joins("JOIN users ON orders.user_id = users.id")
	}

	if needsOrderItemsJoin {
		query = query.Joins("JOIN order_items ON orders.id = order_items.order_id")
	}

	if needsReservationsJoin {
		query = query.Joins("JOIN reservations ON order_items.reservation_id = reservations.id")
	}

	if needsTaxonomiesJoin {
		query = query.Joins("JOIN posts_taxonomies ON posts_taxonomies.post_id = order_items.post_id")
	}

	// Apply time filters on reservation
	if req.StartTime != nil && !req.StartTime.IsZero() {
		query = query.Where("reservations.start_time >= ?", utils.StartOfDayString(*req.StartTime))
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		query = query.Where("reservations.end_time <= ?", utils.EndOfDayString(*req.EndTime))
	}

	// Apply ProductID filter
	if req.ProductID != 0 {
		query = query.Where("CAST(order_items.meta->>'ProductID' AS BIGINT) = ?", req.ProductID)
	}

	// Apply Taxonomies filter
	if len(req.Taxonomies) > 0 {
		query = query.Where("posts_taxonomies.taxonomy_id IN (?)", req.Taxonomies)
	}

	// Apply FullName filter
	if req.FullName != "" {
		query = query.Where(
			"CONCAT(users.first_name, ' ', users.last_name) LIKE ?",
			"%"+strings.TrimSpace(req.FullName)+"%",
		)
	}

	// Group by order ID to avoid duplicates from joins
	if needsOrderItemsJoin {
		query = query.Group("orders.id")
	}

	// Calculate total amount using EXISTS to avoid duplicate amounts from JOINs
	// Apply the same status filter as the main query for consistency
	sumQuery := _i.DB.Main.
		Model(&schema.Order{})

	if req.Status != "" {
		sumQuery = sumQuery.Where("orders.status = ?", req.Status)
	}

	if req.BusinessID > 0 {
		sumQuery = sumQuery.Where(&schema.Order{BusinessID: req.BusinessID})
	}

	if req.UserID > 0 {
		sumQuery = sumQuery.Where(&schema.Order{UserID: req.UserID})
	}

	if req.CouponID > 0 {
		sumQuery = sumQuery.Where(&schema.Order{CouponID: &req.CouponID})
	}

	// Apply HasCoupon filter
	if req.HasCoupon != nil {
		if *req.HasCoupon {
			sumQuery = sumQuery.Where("orders.coupon_id IS NOT NULL")
		} else {
			sumQuery = sumQuery.Where("orders.coupon_id IS NULL")
		}
	}

	// Apply order creation time filters
	if req.OrderStartTime != nil && !req.OrderStartTime.IsZero() {
		sumQuery = sumQuery.Where("orders.created_at >= ?", utils.StartOfDayString(*req.OrderStartTime))
	}

	if req.OrderEndTime != nil && !req.OrderEndTime.IsZero() {
		sumQuery = sumQuery.Where("orders.created_at <= ?", utils.EndOfDayString(*req.OrderEndTime))
	}

	// Apply FullName filter using EXISTS subquery
	if req.FullName != "" {
		sumQuery = sumQuery.Where(
			`EXISTS (
				SELECT 1 FROM users 
				WHERE users.id = orders.user_id 
				AND CONCAT(users.first_name, ' ', users.last_name) LIKE ?
			)`,
			"%"+strings.TrimSpace(req.FullName)+"%",
		)
	}

	// Apply time filters using EXISTS subquery
	if req.StartTime != nil && !req.StartTime.IsZero() {
		sumQuery = sumQuery.Where(
			`EXISTS (
				SELECT 1 FROM order_items 
				JOIN reservations ON order_items.reservation_id = reservations.id 
				WHERE order_items.order_id = orders.id 
				AND reservations.start_time >= ?
			)`,
			utils.StartOfDayString(*req.StartTime),
		)
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		sumQuery = sumQuery.Where(
			`EXISTS (
				SELECT 1 FROM order_items 
				JOIN reservations ON order_items.reservation_id = reservations.id 
				WHERE order_items.order_id = orders.id 
				AND reservations.end_time <= ?
			)`,
			utils.EndOfDayString(*req.EndTime),
		)
	}

	if req.ProductID != 0 {
		sumQuery = sumQuery.Where(
			`EXISTS (
				SELECT 1 FROM order_items 
				WHERE order_items.order_id = orders.id 
				AND CAST(order_items.meta->>'ProductID' AS BIGINT) = ?
			)`,
			req.ProductID,
		)
	} else if len(req.Taxonomies) > 0 {
		sumQuery = sumQuery.Where(
			`EXISTS (
				SELECT 1 FROM order_items 
				JOIN posts_taxonomies ON posts_taxonomies.post_id = order_items.post_id 
				WHERE order_items.order_id = orders.id 
				AND posts_taxonomies.taxonomy_id IN (?)
			)`,
			req.Taxonomies,
		)
	}

	if err = sumQuery.
		Select("COALESCE(CAST(ROUND(SUM(total_amt)) AS BIGINT), 0)").
		Scan(&totalAmount).Error; err != nil {
		return
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query = query.Offset(req.Pagination.Offset)
		query = query.Limit(req.Pagination.Limit)
	}

	if req.BusinessID > 0 {
		query = query.Preload("User")
	}

	err = query.Debug().Preload("Coupon").Preload("OrderItems.Reservation").Order("orders.created_at desc").Find(&orders).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(userID uint64, id uint64) (order *schema.Order, err error) {
	if err := _i.DB.Main.
		Where(&schema.Order{UserID: userID}).
		Preload("OrderItems").
		First(&order, id).
		Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (_i *repo) Create(order *schema.Order, tx *gorm.DB) (orderID uint64, err error) {
	if tx != nil {
		if err := tx.Create(&order).Error; err != nil {
			return 0, err
		}
	} else {
		if err := _i.DB.Main.Create(&order).Error; err != nil {
			return 0, err
		}
	}
	return order.ID, nil
}

func (_i *repo) Update(id uint64, order *schema.Order) (err error) {
	return _i.DB.Main.Model(&schema.Order{}).
		Where(&schema.Order{ID: id, BusinessID: order.BusinessID}).
		Updates(order).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Order{}, id).Error
}

func (_i *repo) BeginTransaction() (*gorm.DB, error) {
	tx := _i.DB.Main.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
