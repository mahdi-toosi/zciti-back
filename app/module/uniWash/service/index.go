package service

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/uniWash/repository"
	"time"
)

type IService interface {
	Reserve(req request.StoreUniWash) (reservationID *uint64, err error)
	IsReservable(req request.StoreUniWash) error
	ValidateReservation(req request.StoreUniWash) error
}

func Service(repo repository.IRepository) IService {
	return &service{
		repo,
	}
}

type service struct {
	Repo repository.IRepository
}

var DefaultSetting = schema.ProductMetaReservation{
	Quantity:  1,
	EndTime:   "07:00:00",
	StartTime: "23:30:00",
	Duration:  (90 * time.Minute) / time.Millisecond,
	Info: map[time.Weekday][]schema.ProductMetaReservationInfoData{
		0: { // sunday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
		1: { // monday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
		2: { // tuesday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
		3: { // wednesday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
		4: { // thursday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
		5: { // friday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
		6: { // saturday
			{From: "07:00:00", To: "08:30:00"},
			{From: "08:30:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:30:00"},
			{From: "11:30:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:30:00"},
			{From: "14:30:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:30:00"},
			{From: "17:30:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:30:00"},
			{From: "20:30:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:30:00"},
		},
	},
}

func (_i *service) Reserve(req request.StoreUniWash) (reservationID *uint64, err error) {
	reservationID, err = _i.Repo.Reserve(req)
	if err != nil {
		return nil, err
	}
	return reservationID, nil
}

func (_i *service) IsReservable(req request.StoreUniWash) (err error) {
	if err = _i.Repo.IsReservable(req); err != nil {
		return err
	}
	return nil
}

func (_i *service) ValidateReservation(req request.StoreUniWash) error {
	invalidErr := &fiber.Error{
		Code:    fiber.StatusBadRequest,
		Message: "invalid reservation",
	}

	parsedDate, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return invalidErr
	}
	day := DefaultSetting.Info[parsedDate.Weekday()]
	for _, hours := range day {
		if hours.From == req.StartTime && hours.To == req.EndTime {
			return nil
		}
	}

	return invalidErr
}

// func (_i *service) Search(req request.Taxonomies, forUser bool) (taxonomies []*response.Taxonomy, paging paginator.Pagination, err error) {
//	results, paging, err := _i.Repo.Search(req)
//	if err != nil {
//		return
//	}
//
//	for _, result := range results {
//		taxonomies = append(taxonomies, response.FromDomain(result, forUser))
//	}
//
//	return
//}
