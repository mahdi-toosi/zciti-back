package cron

import (
	"fmt"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/reservation/repository"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils/config"
	"time"

	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/rs/zerolog"
)

type TurnOffReminderService struct {
	CronSpec   string
	Cfg        *config.Config
	Logger     zerolog.Logger
	Repo       repository.IRepository
	SmsService *internal.MessageWayService // Interface for SMS sending
}

func RunTurnOffMsgReminders(
	cfg *config.Config,
	logger zerolog.Logger,
	Repo repository.IRepository,
	cronService *internal.CronService,
	smsService *internal.MessageWayService,
) *TurnOffReminderService {
	service := &TurnOffReminderService{
		Cfg:        cfg,
		Repo:       Repo,
		Logger:     logger,
		SmsService: smsService,
		CronSpec:   "@every 1m",
	}

	err := cronService.AddJob(service.CronSpec, service.SendTurnOffReminders)
	if err != nil {
		service.Logger.Fatal().Err(err).Msg("failed to add RunTurnOffMsgReminders job")
	}

	return service
}

// SendTurnOffReminders checks and sends reminders for upcoming reservations
func (_s *TurnOffReminderService) SendTurnOffReminders() {
	//_s.Logger.Info().Msg("SendTurnOffReminders")

	// Find reservations within the next 1 hour that haven't been reminded yet
	reservations, err := _s.FindReservationsDueForReminder()
	if err != nil {
		_s.Logger.Err(err).Msg("Failed to fetch reservations for reminders")
		return
	}

	for _, reservation := range reservations {
		// Send reminder and mark as reminded
		if reservation.Product.Post.Status == schema.PostStatusPublished &&
			reservation.Product.Meta.UniWashMachineStatus == schema.UniWashMachineStatusON &&
			reservation.Meta.UniWashLastCommand != "" {
			err := _s.ProcessReservationReminder(*reservation)
			if err != nil {
				_s.Logger.Err(err).Msg("Failed to process reservation reminder")
			}
		}
	}
}

// FindReservationsDueForReminder finds reservations due for reminder
func (_s *TurnOffReminderService) FindReservationsDueForReminder() ([]*schema.Reservation, error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	t := time.Now().In(loc)
	//t := time.Date(2025, 8, 9, 10, 40, 5, 0, loc)

	endTime := t.Add(20 * time.Minute).Truncate(time.Minute)
	startTime := endTime.Add(-1 * time.Hour)
	reminderNotSent := false

	reservations, _, err := _s.Repo.GetAll(request.Reservations{
		EndTime:             &endTime,
		StartTime:           &startTime,
		Status:              schema.ReservationStatusReserved,
		TurnOffReminderSent: &reminderNotSent,
	})

	//_s.Logger.Info().Msgf("turn off reservations len %d", len(reservations))

	if err != nil {
		return nil, err
	}

	return reservations, err
}

// ProcessReservationReminder sends SMS and updates reminder status
func (_s *TurnOffReminderService) ProcessReservationReminder(reservation schema.Reservation) error {
	// Send reminder message first
	_, err := _s.SmsService.Send(MessageWay.Message{
		Provider:   5, // با سرشماره 5000
		TemplateID: 16621,
		Method:     "sms",
		Mobile:     fmt.Sprintf("0%d", reservation.User.Mobile),
	})
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// Mark as sent only AFTER successful SMS delivery
	if err := _s.Repo.MarkTurnOffReminderSent(reservation.ID); err != nil {
		_s.Logger.Err(err).Uint64("reservationID", reservation.ID).Msg("SMS sent but failed to mark reminder as sent")
		return fmt.Errorf("failed to mark reminder as sent: %w", err)
	}

	return nil
}
