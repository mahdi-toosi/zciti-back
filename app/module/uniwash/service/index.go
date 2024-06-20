package service

import (
	"go-fiber-starter/app/database/schema"
	oirequest "go-fiber-starter/app/module/orderItem/request"
	prepository "go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/uniwash/repository"
	"go-fiber-starter/app/module/uniwash/request"
	"go-fiber-starter/app/module/uniwash/response"
	"go-fiber-starter/utils/paginator"
	"time"

	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/gofiber/fiber/v2"
)

type IService interface {
	ValidateReservation(req oirequest.OrderItem) error
	IsReservable(req oirequest.OrderItem, businessID uint64) error
	ReserveReservation(req oirequest.OrderItem, userID uint64, businessID uint64) (reservationID *uint64, err error)
	Reserve(reservationID uint64) error
	SendCommand(req request.SendCommand, isForUser bool) error
	IndexReservedMachines(req request.ReservedMachinesRequest) (reserved []*response.Reservation, paging paginator.Pagination, err error)
	CheckLastCommandStatus(businessID uint64, reservationID uint64) (status *MessageWay.StatusResponse, err error)
}

func Service(repo repository.IRepository, productRepo prepository.IRepository, messageWay *MessageWay.App) IService {
	return &service{repo, messageWay, productRepo}
}

type service struct {
	Repo        repository.IRepository
	MessageWay  *MessageWay.App
	ProductRepo prepository.IRepository
}

func (_i *service) ReserveReservation(req oirequest.OrderItem, userID uint64, businessID uint64) (reservationID *uint64, err error) {
	reservationID, err = _i.Repo.ReserveReservation(req, userID, businessID)
	if err != nil {
		return nil, err
	}
	return reservationID, nil
}

func (_i *service) IsReservable(req oirequest.OrderItem, businessID uint64) (err error) {
	if err = _i.Repo.IsReservable(req, businessID); err != nil {
		return err
	}
	return nil
}

func (_i *service) ValidateReservation(req oirequest.OrderItem) error {
	invalidErr := &fiber.Error{
		Code:    fiber.StatusBadRequest,
		Message: "در این بازه زمانی شما اجازه رزرو ندارید",
	}

	parsedDate, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return invalidErr
	}

	if req.GetStartDateTime().Before(time.Now()) {
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
		if !time.Now().After(reservation.StartTime.Add(-10 * time.Minute)) {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "در بازه زمانی که رزرو کرده اید، دوباره تلاش کنید",
			}
		}

		if !time.Now().Before(reservation.StartTime.Add(10 * time.Minute)) {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "شما تا نهایاتا 10 دقیقه پس از شروع زمان فرصت داشتید به دستگاه فرمان بدهید.",
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
	if err != nil || product.Meta.UniWashMachineStatus == schema.UniWashMachineStatusOFF {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "در حال حاضر دستگاه دچار مشکل شده و در دسترس نیست، در صورتی که دستگاه را رزرو کرده اید، با پشتیبانی تماس بگیرید",
		}
	}

	send, err := _i.MessageWay.Send(MessageWay.Message{
		Provider:   3, // با سرشماره 9000
		TemplateID: 8698,
		Method:     "sms",
		Params:     []string{commandProxy[req.Command]},
		Mobile:     product.Meta.UniWashMobileNumber,
	})
	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "ارسال دستور با خطا مواجه شد، دوباره امتحان کنید."}
	}

	// if isForUser {
	t := time.Now().UTC()
	reservation.Meta.UniWashLastCommandTime = &t
	reservation.Meta.UniWashLastCommand = req.Command
	reservation.Meta.UniWashLastCommandReferenceID = send.ReferenceID
	if err := _i.Repo.UpdateReservation(reservation); err != nil {
		return err
	}
	//}

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

func (_i *service) Reserve(reservationID uint64) (err error) {
	if err := _i.Repo.Reserve(reservationID); err != nil {
		return err
	}
	return nil
}

func (_i *service) CheckLastCommandStatus(businessID uint64, reservationID uint64) (status *MessageWay.StatusResponse, err error) {
	reservation, err := _i.Repo.GetSingleReservation(businessID, reservationID)
	if err != nil {
		return nil, err
	}

	if reservation.Meta.UniWashLastCommand == "" {
		return &MessageWay.StatusResponse{
			OTPStatus: "دستوری ارسال نشده",
			Status:    "danger",
		}, nil
	}

	return _i.MessageWay.GetStatus(MessageWay.StatusRequest{ReferenceID: reservation.Meta.UniWashLastCommandReferenceID})
}

var commandProxy = map[schema.UniWashCommand]string{
	schema.UniWashCommandON:        "on",
	schema.UniWashCommandOFF:       "off",
	schema.UniWashCommandMoreWater: "7",
}

var DefaultSetting = schema.ProductMetaReservation{
	Quantity:  1,
	EndTime:   "00:00:00",
	StartTime: "00:00:00",
	Duration:  (60 * time.Minute) / time.Millisecond,
	Info: map[time.Weekday][]schema.ProductMetaReservationInfoData{
		0: { // sunday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
		1: { // monday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
		2: { // tuesday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
		3: { // wednesday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
		4: { // thursday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
		5: { // friday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
		6: { // saturday
			{From: "00:00:00", To: "01:00:00"},
			{From: "01:00:00", To: "02:00:00"},
			{From: "02:00:00", To: "03:00:00"},
			{From: "03:00:00", To: "04:00:00"},
			{From: "04:00:00", To: "05:00:00"},
			{From: "05:00:00", To: "06:00:00"},
			{From: "06:00:00", To: "07:00:00"},
			{From: "07:00:00", To: "08:00:00"},
			{From: "08:00:00", To: "09:00:00"},
			{From: "09:00:00", To: "10:00:00"},
			{From: "10:00:00", To: "11:00:00"},
			{From: "11:00:00", To: "12:00:00"},
			{From: "12:00:00", To: "13:00:00"},
			{From: "13:00:00", To: "14:00:00"},
			{From: "14:00:00", To: "15:00:00"},
			{From: "15:00:00", To: "16:00:00"},
			{From: "16:00:00", To: "17:00:00"},
			{From: "17:00:00", To: "18:00:00"},
			{From: "18:00:00", To: "19:00:00"},
			{From: "19:00:00", To: "20:00:00"},
			{From: "20:00:00", To: "21:00:00"},
			{From: "21:00:00", To: "22:00:00"},
			{From: "22:00:00", To: "23:00:00"},
			{From: "23:00:00", To: "00:00:00"},
		},
	},
}
