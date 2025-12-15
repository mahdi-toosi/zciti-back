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

type TurnOnReminderService struct {
	CronSpec   string
	Cfg        *config.Config
	Logger     zerolog.Logger
	Repo       repository.IRepository
	SmsService *internal.MessageWayService // Interface for SMS sending
}

func RunTurnOnMsgReminders(
	cfg *config.Config,
	logger zerolog.Logger,
	Repo repository.IRepository,
	cronService *internal.CronService,
	smsService *internal.MessageWayService,
) *TurnOnReminderService {
	service := &TurnOnReminderService{
		Cfg:        cfg,
		Repo:       Repo,
		Logger:     logger,
		SmsService: smsService,
		CronSpec:   "@every 1m",
	}

	err := cronService.AddJob(service.CronSpec, service.SendTurnOnReminders)
	if err != nil {
		service.Logger.Fatal().Err(err).Msg("failed to add RunTurnOnMsgReminders job")
	}

	return service
}

// SendTurnOnReminders checks and sends reminders for upcoming reservations
func (_s *TurnOnReminderService) SendTurnOnReminders() {
	//_s.Logger.Info().Msg("SendTurnOnReminders")

	//Find reservations within the next 1 hour that haven't been reminded yet
	reservations, err := _s.FindReservationsDueForReminder()
	if err != nil {
		_s.Logger.Err(err).Msg("Failed to fetch reservations for reminders")
		return
	}

	for _, reservation := range reservations {
		//_s.Logger.Info().Msgf("reservation id => %d", reservation.ID)
		// Send reminder and mark as reminded
		if reservation.Product.Post.Status == schema.PostStatusPublished &&
			reservation.Product.Meta.UniWashMachineStatus == schema.UniWashMachineStatusON {
			err := _s.ProcessReservationReminder(*reservation)
			if err != nil {
				_s.Logger.Err(err).Msg("Failed to process reservation reminder")
			}
		}
	}
}

// FindReservationsDueForReminder finds reservations due for reminder
func (_s *TurnOnReminderService) FindReservationsDueForReminder() ([]*schema.Reservation, error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	t := time.Now().In(loc)
	//t := time.Date(2025, 9, 8, 9, 26, 5, 0, loc)

	startTime := t.Add(60 * time.Minute).Truncate(time.Minute)
	endTime := startTime.Add(time.Hour)
	reminderNotSent := false

	reservations, _, err := _s.Repo.GetAll(request.Reservations{
		EndTime:            &endTime,
		StartTime:          &startTime,
		Status:             schema.ReservationStatusReserved,
		TurnOnReminderSent: &reminderNotSent,
	})
	if err != nil {
		return nil, err
	}

	//_s.Logger.Info().Msgf("turn on reservations len %d", len(reservations))

	return reservations, err
}

// ProcessReservationReminder sends SMS and updates reminder status
func (_s *TurnOnReminderService) ProcessReservationReminder(reservation schema.Reservation) error {
	_s.Logger.Info().Msg("ProcessReservationReminder")

	// Send reminder message first
	_, err := _s.SmsService.Send(MessageWay.Message{
		Provider:   5, // با سرشماره 5000
		TemplateID: 16620,
		Method:     "sms",
		Mobile:     fmt.Sprintf("0%d", reservation.User.Mobile),
	})
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	// Mark as sent only AFTER successful SMS delivery
	if err := _s.Repo.MarkTurnOnReminderSent(reservation.ID); err != nil {
		_s.Logger.Err(err).Uint64("reservationID", reservation.ID).Msg("SMS sent but failed to mark reminder as sent")
		return fmt.Errorf("failed to mark reminder as sent: %w", err)
	}

	return nil
}
