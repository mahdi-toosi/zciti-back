package repository

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetReservation(req request.SendCommand) (reservation *schema.Reservation, err error)
	UpdateReservation(reservation *schema.Reservation) error
	Reserve(req request.StoreUniWash) (reservationID *uint64, err error)
	IsReservable(req request.StoreUniWash) error
	IndexReservedMachines(req request.ReservedMachinesRequest) (reservations []*schema.Reservation, paging paginator.Pagination, err error)
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
		}).
		First(&reservation).Error; err != nil {
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

func (_i *repo) Reserve(req request.StoreUniWash) (reservationID *uint64, err error) {
	r := schema.Reservation{
		UserID:     req.UserID,
		ProductID:  req.ProductID,
		BusinessID: req.BusinessID,
		EndTime:    req.GetEndDateTime(),
		StartTime:  req.GetStartDateTime(),
		Status:     schema.ReservationStatusReserved,
	}
	if err = _i.DB.Main.Create(&r).Error; err != nil {
		return nil, err
	}

	return &r.ID, nil
}

func (_i *repo) IsReservable(req request.StoreUniWash) error {
	var reservation schema.Reservation
	if err := _i.DB.Main.
		Where(&schema.Reservation{
			ProductID:  req.ProductID,
			BusinessID: req.BusinessID,
			EndTime:    req.GetEndDateTime(),
			StartTime:  req.GetStartDateTime(),
			Status:     schema.ReservationStatusReserved,
		}).
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
		Where(&schema.Reservation{BusinessID: req.BusinessID, UserID: req.UserID})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.
		Preload("Product").
		Order("created_at desc").Find(&reservations).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}
