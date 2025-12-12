package repository

import (
	"go-fiber-starter/app/database/schema"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IRepository interface {
	GetReservation(req request.SendCommand) (reservation *schema.Reservation, err error)
	GetSingleReservation(BusinessID uint64, id uint64) (reservation *schema.Reservation, err error)
	UpdateReservation(reservation *schema.Reservation) error
	ReserveReservation(req oirequest.OrderItem, userID uint64, businessID uint64) (reservationID *uint64, err error)
	IsReservable(req oirequest.OrderItem, businessID uint64) error
	IndexReservedMachines(req request.ReservedMachinesRequest) (reservations []*schema.Reservation, paging paginator.Pagination, err error)
	Reserve(reservationID uint64) error
}

func Repository(db *database.Database) IRepository {
	return &repo{db}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetReservation(req request.SendCommand) (reservation *schema.Reservation, err error) {
	if err := _i.DB.Main.
		Where(&schema.Reservation{
			UserID:     req.UserID,
			ProductID:  req.ProductID,
			BusinessID: req.BusinessID,
			ID:         req.ReservationID,
			Status:     schema.ReservationStatusReserved,
		}).Debug().
		First(&reservation).Error; err != nil {
		return nil, err
	}

	return reservation, nil
}

func (_i *repo) GetSingleReservation(BusinessID uint64, id uint64) (reservation *schema.Reservation, err error) {
	if err := _i.DB.Main.
		Where(&schema.Reservation{
			BusinessID: BusinessID,
		}).
		Preload("User").
		Preload("Product").
		First(&reservation, id).Error; err != nil {
		return nil, err
	}

	return reservation, nil
}

func (_i *repo) UpdateReservation(reservation *schema.Reservation) (err error) {
	if err := _i.DB.Main.Model(&schema.Reservation{}).
		Where(&schema.Reservation{ID: reservation.ID, BusinessID: reservation.BusinessID}).
		Updates(reservation).Error; err != nil {
		return err
	}
	return nil
}

func (_i *repo) ReserveReservation(req oirequest.OrderItem, userID uint64, businessID uint64) (reservationID *uint64, err error) {
	r := schema.Reservation{
		UserID:     userID,
		BusinessID: businessID,
		ProductID:  req.ProductID,
		EndTime:    req.GetEndDateTime(),
		StartTime:  req.GetStartDateTime(),
		Status:     schema.ReservationStatusReserved,
		Base:       schema.Base{DeletedAt: gorm.DeletedAt{Time: time.Now().Add(10 * time.Minute), Valid: true}},
	}
	if err = _i.DB.Main.Create(&r).Error; err != nil {
		return nil, err
	}

	return &r.ID, nil
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

func (_i *repo) IndexReservedMachines(req request.ReservedMachinesRequest) (reservations []*schema.Reservation, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Reservation{}).
		Where(&schema.Reservation{BusinessID: req.BusinessID})

	if req.ProductID > 0 {
		query.Where(&schema.Reservation{ProductID: req.ProductID})
	}

	if !req.Date.IsZero() {
		startTime := utils.StartOfDay(req.Date)
		endTime := utils.EndOfDay(req.Date)
		query.Where("start_time BETWEEN ? AND ?", startTime, endTime)
	}

	if req.With == "reservedReservations" {
		query.Unscoped().Where("deleted_at > ? OR deleted_at IS NULL", time.Now())
	} else if req.UserID > 0 {
		query.Where(&schema.Reservation{UserID: req.UserID}).
			Preload("Product", func(db *gorm.DB) *gorm.DB {
				return db.Unscoped() // This will include soft-deleted posts
			}).
			Preload("Product.Post", func(db *gorm.DB) *gorm.DB {
				return db.Unscoped() // This will include soft-deleted posts
			}).
			Order("start_time desc")
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	if err = query.Find(&reservations).Error; err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) Reserve(id uint64) (err error) {
	if err := _i.DB.Main.Model(&schema.Reservation{}).Unscoped().
		Where(&schema.Reservation{ID: id}).
		Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}
