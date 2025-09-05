package repository

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/app/module/reservation/response"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
	"strings"
	"time"
)

type IRepository interface {
	GetAll(req request.Reservations) (reservations []*schema.Reservation, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (reservation *response.Reservation, err error)
	Create(reservation *schema.Reservation) (err error)
	Update(id uint64, reservation *schema.Reservation) (err error)
	Delete(id uint64) (err error)
	IsReservable(req oirequest.OrderItem, businessID uint64) error
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Reservations) (reservations []*schema.Reservation, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&schema.Reservation{})

	if req.WithUsageCount == 1 {
		// || req.UsageCount > 0
		query.Select(`reservations.*, COUNT(*) OVER (PARTITION BY reservations.user_id) as user_usage_count`)
	}

	//if req.UsageCount > 0 {
	//	// Use a subquery to filter by user count
	//	userSubquery := _i.DB.Main.Model(&schema.Reservation{}).
	//		Select("user_id").
	//		Group("user_id").
	//		Having("COUNT(*) = ?", req.UsageCount)
	//
	//	query.Where("reservations.user_id IN (?)", userSubquery)
	//}

	if req.ProductID != 0 {
		query.Where(&schema.Reservation{ProductID: req.ProductID})
	}

	if req.BusinessID != 0 {
		query.Where(&schema.Reservation{BusinessID: req.BusinessID})
	}

	if req.Status != "" {
		query.Where(&schema.Reservation{Status: req.Status})
	}

	if req.FullName != "" || req.Mobile != "" {
		query.Joins("JOIN users ON reservations.user_id = users.id")
		if req.FullName != "" {
			query.Where(
				"concat(users.first_name, ' ', users.last_name) LIKE ?",
				strings.ReplaceAll("%"+req.FullName+"%", " ", "%"),
			)
		}
		if req.Mobile != "" {
			query.Where("users.mobile = ?", req.Mobile)
		}
	}

	if req.StartTime != nil && !req.StartTime.IsZero() {
		query.Where("start_time >= ?", req.StartTime)
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		query.Where("end_time <= ?", req.EndTime)
	}

	// Handle taxonomy filtering
	var taxonomyConditions []uint64
	if req.CityID > 0 {
		taxonomyConditions = append(taxonomyConditions, req.CityID)
	}
	if req.WorkspaceID > 0 {
		taxonomyConditions = append(taxonomyConditions, req.WorkspaceID)
	}
	if req.DormitoryID > 0 {
		taxonomyConditions = append(taxonomyConditions, req.DormitoryID)
	}

	//if len(taxonomyConditions) > 0 {
	//	// Join through: Reservation -> Product -> Post -> Taxonomy
	//	for i, taxonomyID := range taxonomyConditions {
	//		alias := fmt.Sprintf("pt%d", i)
	//		query.Joins(fmt.Sprintf(`
	//			JOIN products p%d ON p%d.id = reservations.product_id
	//			JOIN posts post%d ON post%d.id = p%d.post_id
	//			JOIN posts_taxonomies %s ON %s.post_id = post%d.id AND %s.taxonomy_id = ?`,
	//			i, i, i, i, i, alias, alias, i, alias), taxonomyID)
	//	}
	//
	//	// Preload the full taxonomy path for the response
	//	query.Preload("Product.Post.Taxonomies", "id IN ?", taxonomyConditions)
	//}

	if len(taxonomyConditions) > 0 {
		// Use EXISTS with correlated subquery
		existsSubquery := _i.DB.Main.Model(&schema.Post{}).
			Select("1").
			Joins("JOIN posts_taxonomies pt ON pt.post_id = posts.id").
			Where("posts.id = products.post_id").
			Where("pt.taxonomy_id IN ?", taxonomyConditions).
			Group("posts.id").
			Having("COUNT(DISTINCT pt.taxonomy_id) = ?", len(taxonomyConditions))

		productExistsSubquery := _i.DB.Main.Model(&schema.Product{}).
			Select("1").
			Where("products.id = reservations.product_id").
			Where("EXISTS (?)", existsSubquery)

		query.Where("EXISTS (?)", productExistsSubquery)
		query.Preload("Product.Post.Taxonomies", "id IN ?", taxonomyConditions)
	}

	if req.Pagination != nil && req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.
		Preload("User").
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // This will include soft-deleted products
		}).
		Preload("Product.Post", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // This will include soft-deleted posts
		}).
		Order("start_time desc").Find(&reservations).Error
	if err != nil {
		return
	}

	if req.Pagination != nil {
		paging = *req.Pagination
	}

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (reservation *response.Reservation, err error) {
	if err := _i.DB.Main.
		Where(&schema.Reservation{BusinessID: businessID}).
		First(&reservation, id).
		Error; err != nil {
		return nil, err
	}

	return reservation, nil
}

func (_i *repo) Create(reservation *schema.Reservation) (err error) {
	return _i.DB.Main.Create(reservation).Error
}

func (_i *repo) Update(id uint64, reservation *schema.Reservation) (err error) {
	return _i.DB.Main.Model(&schema.Reservation{}).
		Where(&schema.Reservation{ID: id, BusinessID: reservation.BusinessID}).
		Updates(reservation).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Reservation{}, id).Error
}

func (_i *repo) IsReservable(req oirequest.OrderItem, businessID uint64) error {
	var reservation schema.Reservation
	if err := _i.DB.Main.
		Where(&schema.Reservation{
			BusinessID: businessID,
			ProductID:  req.ProductID,
			EndTime:    req.GetEndDateTime(),
			StartTime:  req.GetStartDateTime(),
			Status:     schema.ReservationStatusReserved,
		}).
		Unscoped().Where("deleted_at > ? OR deleted_at IS NULL", time.Now()).
		First(&reservation).Error; err == nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "این ساعت دستگاه رزرو شده است",
		}
	}

	return nil
}
