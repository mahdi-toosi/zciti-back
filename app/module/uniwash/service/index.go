package service

import (
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	prepository "go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/uniwash/repository"
	"go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/app/module/uniwash/response"
	"go-fiber-starter/utils/paginator"
	"time"
)

type IService interface {
	Reserve(req request.StoreUniWash) (reservationID *uint64, err error)
	IsReservable(req request.StoreUniWash) error
	ValidateReservation(req request.StoreUniWash) error
	SendCommand(req request.SendCommand, isForUser bool) error
	IndexReservedMachines(req request.ReservedMachinesRequest) (reserved []*response.Reservation, paging paginator.Pagination, err error)
}

func Service(repo repository.IRepository, productRepo prepository.IRepository, messageWay *MessageWay.App) IService {
	return &service{repo, messageWay, productRepo}
}

type service struct {
	Repo        repository.IRepository
	MessageWay  *MessageWay.App
	ProductRepo prepository.IRepository
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

func (_i *service) SendCommand(req request.SendCommand, isForUser bool) (err error) {
	var reservation *schema.Reservation

	if isForUser {
		reservation, err = _i.Repo.GetReservation(req)
		if err != nil {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "شما این دستگاه را رزرو نکرده اید",
			}
		}
		// TODO talk about this condition with esmaiil and hosein
		// check the 10 min before start time is after now
		if !time.Now().After(reservation.StartTime.Add(-10*time.Minute)) ||
			!time.Now().Before(reservation.EndTime.Add(-10*time.Minute)) {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "در بازه زمانی که رزرو کرده اید، دوباره تلاش کنید",
			}
		}

		if reservation.Meta.UniWashLastCommand == schema.UniWashCommandON {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "دستگاه در حال حاضر روشن می باشد",
			}
		}
	}

	product, err := _i.ProductRepo.GetOneVariant(req.BusinessID, req.ProductID)
	if err != nil || product.Meta.UniWashMachineStatus == schema.UniWashCommandOffline {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "در حال حاضر رزرو این ماشین ممکن نیست",
		}
	}

	send, err := _i.MessageWay.Send(MessageWay.Message{
		TemplateID: 8698,
		Method:     "sms",
		Params:     []string{commandProxy[schema.UniWashCommandON]},
		Mobile:     product.Meta.UniWashMobileNumber,
	})
	if err != nil {
		return err
	}

	if isForUser {
		reservation.Meta.UniWashLastCommandTime = time.Now().UTC()
		reservation.Meta.UniWashLastCommand = schema.UniWashCommandON
		reservation.Meta.UniWashLastCommandReferenceID = send.ReferenceID
		if err := _i.Repo.UpdateReservation(reservation); err != nil {
			return err
		}
	}

	product.Meta.UniWashMachineStatus = req.Command
	if err := _i.ProductRepo.Update(product); err != nil {
		return err
	}

	return nil
}

func (_i *service) IndexReservedMachines(req request.ReservedMachinesRequest) (reserved []*response.Reservation, paging paginator.Pagination, err error) {

	results, paging, err := _i.Repo.IndexReservedMachines(req)
	if err != nil {
		return
	}

	for _, result := range results {
		reserved = append(reserved, response.FromDomain(result))
	}

	return
}

// TODO fill the correct command
var commandProxy = map[schema.UniWashCommand]string{
	schema.UniWashCommandON:        "7",
	schema.UniWashCommandOFF:       "2",
	schema.UniWashCommandMoreWater: "3",
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
