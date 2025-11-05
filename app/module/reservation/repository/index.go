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

	if req.CityID > 0 || req.WorkspaceID > 0 || req.DormitoryID > 0 {
		var filters []uint64
		if req.DormitoryID > 0 {
			filters = append(filters, req.DormitoryID)
		} else if req.WorkspaceID > 0 {
			filters = append(filters, req.WorkspaceID)
		} else if req.CityID > 0 {
			filters = append(filters, req.CityID)
		}

		query.Joins("JOIN products ON reservations.product_id = products.id")
		query.Joins("JOIN posts_taxonomies ON posts_taxonomies.post_id = products.post_id")
		query.Where("posts_taxonomies.taxonomy_id IN (?)", filters)
		query.Preload("Product.Post.Taxonomies")
		query.Group("reservations.id")
	}

	if len(req.Posts) > 0 {
		if req.CityID == 0 && req.WorkspaceID == 0 && req.DormitoryID == 0 {
			query.Joins("JOIN products ON reservations.product_id = products.id")
			query.Group("reservations.id")
		}
		query.Where("products.post_id IN (?)", req.Posts)
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
