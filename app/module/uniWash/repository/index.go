package repository

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/internal/bootstrap/database"
)

type IRepository interface {
	Reserve(req request.StoreUniWash) (reservationID *uint64, err error)
	IsReservable(req request.StoreUniWash) error
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
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
		First(&reservation).
		Error; err == nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "already reserved",
		}
	}

	return nil
}
