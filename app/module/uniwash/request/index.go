package request

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"time"
)

type SendCommand struct {
	BusinessID    uint64
	UserID        uint64
	ReservationID uint64                `example:"1" validate:"required,number"`
	ProductID     uint64                `example:"1" validate:"required,number"`
	Command       schema.UniWashCommand `example:"ON" validate:"required,oneof=ON OFF REWASH EVACUATION"`
}

type ReservedMachinesRequest struct {
	BusinessID uint64
	UserID     uint64
	ProductID  uint64
	Date       time.Time
	With       string
	Pagination *paginator.Pagination
}

type StoreUniWash struct {
	Date       string
	UserID     uint64
	PostID     uint64
	EndTime    string
	StartTime  string
	ProductID  uint64
	BusinessID uint64
}

func (s StoreUniWash) GetStartDateTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Tehran")
	startTime, _ := time.ParseInLocation(time.DateTime, s.Date+" "+s.StartTime, loc)
	return startTime
}

func (s StoreUniWash) GetEndDateTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Tehran")
	endTime, _ := time.ParseInLocation(time.DateTime, s.Date+" "+s.EndTime, loc)

	// if hour is 00 should store in next day
	if endTime.Hour() == 0 {
		endTime = endTime.Add(24 * time.Hour)
	}

	return endTime
}
