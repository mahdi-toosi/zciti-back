package service

import (
	"fmt"
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
	GetReservationOptions() (reservationOptions schema.ProductMetaReservationOptions)
	SendDeviceIsOffMsgToUser(businessID uint64, reservationID uint64) (err error)
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

	loc, _ := time.LoadLocation("Asia/Tehran")
	parsedDate, err := time.ParseInLocation(time.DateOnly, req.Date, loc)
	if err != nil {
		return invalidErr
	}

	if req.GetStartDateTime().Before(time.Now()) {
		return invalidErr
	}

	reservationOptions := _i.GetReservationOptions()
	day := reservationOptions[parsedDate.Weekday()]
	for _, hours := range day {
		if hours.From == req.StartTime && hours.To == req.EndTime {
			return nil
		}
	}

	return invalidErr
}

func (_i *service) SendCommand(req request.SendCommand, isForUser bool) (err error) {
	reservation, err := _i.Repo.GetReservation(req)

	if isForUser {
		if err != nil {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "شما این دستگاه را رزرو نکرده اید",
			}
		}

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
				Message: "شما تا نهایتاً 10 دقیقه پس از شروع زمان فرصت داشتید به دستگاه فرمان بدهید.",
			}
		}

		if reservation.Meta.UniWashLastCommand == schema.UniWashCommandON &&
			!time.Now().After(reservation.Meta.UniWashLastCommandTime.Add(30*time.Second)) {
			return &fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "در صورت روشن نشدن دستگاه ۳۰ ثانیه بعد دوباره درخواست دهید.",
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
		Provider:   3, // با سر شماره 9000
		TemplateID: 8698,
		Method:     "sms",
		Params:     []string{commandProxy[req.Command]},
		Mobile:     product.Meta.UniWashMobileNumber,
	})

	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "ارسال دستور با خطا مواجه شد، دوباره امتحان کنید."}
	}

	if send.Status == "error" {
		return &fiber.Error{Code: fiber.StatusServiceUnavailable, Message: "ارسال دستور با خطا مواجه شد کد ۵۰۶، با پشتیبانی در میان بگذارید."}
	}

	t := time.Now().UTC()
	reservation.Meta.UniWashLastCommandTime = &t
	reservation.Meta.UniWashLastCommand = req.Command
	reservation.Meta.UniWashLastCommandReferenceID = send.ReferenceID
	if err := _i.Repo.UpdateReservation(reservation); err != nil {
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

func (_i *service) SendDeviceIsOffMsgToUser(businessID uint64, reservationID uint64) (err error) {
	reservation, err := _i.Repo.GetSingleReservation(businessID, reservationID)
	if err != nil {
		return err
	}

	send, err := _i.MessageWay.Send(MessageWay.Message{
		Provider:   5, // با سرشماره 5000
		TemplateID: 16622,
		Method:     "sms",
		Params:     []string{reservation.Product.Meta.SKU},
		Mobile:     fmt.Sprintf("0%d", reservation.User.Mobile),
	})

	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "ارسال دستور با خطا مواجه شد، دوباره امتحان کنید."}
	}

	if send.Status == "error" {
		return &fiber.Error{Code: fiber.StatusServiceUnavailable, Message: "ارسال دستور با خطا مواجه شد کد ۵۰۶، با پشتیبانی در میان بگذارید."}
	}

	return nil
}

func (_i *service) GetReservationOptions() (reservationOptions schema.ProductMetaReservationOptions) {
	options := schema.ProductMetaReservationOptions{}
	// Sunday is 0
	dayNumbers := []time.Weekday{0, 1, 2, 3, 4, 5, 6}

	for _, dayNum := range dayNumbers {
		options[dayNum] = []schema.ProductMetaDeviceHour{}
		for hour := 0; hour < 24; hour++ {
			to := ""
			if hour+1 == 24 {
				to = "00:00:00"
			} else {
				to = fmt.Sprintf("%02d:00:00", hour+1)
			}

			deviceHour := schema.ProductMetaDeviceHour{
				ID:   fmt.Sprintf("%d-%02d", dayNum, hour),
				From: fmt.Sprintf("%02d:00:00", hour),
				To:   to,
			}
			options[dayNum] = append(options[dayNum], deviceHour)
		}

	}
	return options
}

var commandProxy = map[schema.UniWashCommand]string{
	schema.UniWashCommandON:         "on",
	schema.UniWashCommandOFF:        "off",
	schema.UniWashCommandRewash:     "10",
	schema.UniWashCommandEvacuation: "9",
}
