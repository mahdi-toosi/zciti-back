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
	cronSpec   string
	cfg        *config.Config
	logger     zerolog.Logger
	repo       repository.IRepository
	smsService *internal.MessageWayService // Interface for SMS sending
}

func RunTurnOnMsgReminders(
	cfg *config.Config,
	logger zerolog.Logger,
	Repo repository.IRepository,
	cronService *internal.CronService,
	smsService *internal.MessageWayService,
) *TurnOnReminderService {
	service := &TurnOnReminderService{
		cfg:        cfg,
		repo:       Repo,
		logger:     logger,
		smsService: smsService,
		cronSpec:   "@every 1m",
	}

	err := cronService.AddJob(service.cronSpec, service.SendTurnOnReminders)
	if err != nil {
		service.logger.Fatal().Err(err).Msg("failed to add RunTurnOnMsgReminders job")
	}

	return service
}

// SendTurnOnReminders checks and sends reminders for upcoming reservations
func (_s *TurnOnReminderService) SendTurnOnReminders() {
	//_s.logger.Info().Msg("SendTurnOnReminders")

	//Find reservations within the next 1 hour that haven't been reminded yet
	reservations, err := _s.findReservationsDueForReminder()
	if err != nil {
		_s.logger.Err(err).Msg("Failed to fetch reservations for reminders")
		return
	}

	for _, reservation := range reservations {
		//_s.logger.Info().Msgf("reservation id => %d", reservation.ID)
		// Send reminder and mark as reminded
		if reservation.Product.Post.Status == schema.PostStatusPublished &&
			reservation.Product.Meta.UniWashMachineStatus == schema.UniWashMachineStatusON {
			err := _s.processReservationReminder(*reservation)
			if err != nil {
				_s.logger.Err(err).Msg("Failed to process reservation reminder")
			}
		}
	}
}

// findReservationsDueForReminder finds reservations due for reminder
func (_s *TurnOnReminderService) findReservationsDueForReminder() ([]*schema.Reservation, error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	t := time.Now().In(loc)
	//t := time.Date(2025, 9, 8, 9, 26, 5, 0, loc)

	startTime := t.Add(60 * time.Minute).Truncate(time.Minute)
	endTime := startTime.Add(time.Hour)

	reservations, _, err := _s.repo.GetAll(request.Reservations{
		EndTime:   &endTime,
		StartTime: &startTime,
		Status:    schema.ReservationStatusReserved,
	})
	if err != nil {
		return nil, err
	}

	//s.logger.Info().Msgf("turn on reservations len %d", len(reservations))

	return reservations, err
}

// processReservationReminder sends SMS and updates reminder status
func (_s *TurnOnReminderService) processReservationReminder(reservation schema.Reservation) error {
	_s.logger.Info().Msg("processReservationReminder")
	// Prepare reminder message
	_, err := _s.smsService.Send(MessageWay.Message{
		Provider:   5, // با سرشماره 5000
		TemplateID: 16620,
		Method:     "sms",
		Mobile:     fmt.Sprintf("0%d", reservation.User.Mobile),
	})
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}
