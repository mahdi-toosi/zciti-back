package repository

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/app/module/reservation/response"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
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

	if req.ProductID != 0 {
		query.Where(&schema.Reservation{ProductID: req.ProductID})
	}

	if req.BusinessID != 0 {
		query.Where(&schema.Reservation{BusinessID: req.BusinessID})
	}

	if req.FullName != "" || req.Mobile != "" {
		query.Joins("JOIN users ON reservations.user_id = users.id")
		if req.FullName != "" {
			query.Where("concat(users.first_name, ' ', users.last_name) LIKE ?", "%"+req.FullName+"%")
		}
		if req.Mobile != "" {
			query.Where("users.mobile = ?", req.Mobile)
		}
	}

	utils.Log(req.StartTime.IsZero())
	if !req.StartTime.IsZero() {
		query.Where("start_time::Date = ?", req.StartTime)
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.
		Preload("User").Preload("Product.Post").Order("start_time desc").Find(&reservations).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

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
