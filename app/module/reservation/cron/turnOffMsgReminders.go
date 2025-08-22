package cron

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/rs/zerolog"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/reservation/repository"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils/config"
	"time"
)

type TurnOffReminderService struct {
	cronSpec   string
	cfg        *config.Config
	logger     zerolog.Logger
	smsService *MessageWay.App // Interface for SMS sending
	repo       repository.IRepository
}

func RunTurnOffMsgReminders(
	cfg *config.Config,
	logger zerolog.Logger,
	Repo repository.IRepository,
	smsService *MessageWay.App,
	cronService *internal.CronService,
) *TurnOffReminderService {
	service := &TurnOffReminderService{
		cfg:        cfg,
		repo:       Repo,
		logger:     logger,
		smsService: smsService,
		cronSpec:   "@every 1m",
	}

	err := cronService.AddJob(service.cronSpec, service.SendTurnOffReminders)
	if err != nil {
		service.logger.Fatal().Err(err).Msg("failed to add RunTurnOffMsgReminders job")
	}

	return service
}

// SendTurnOffReminders checks and sends reminders for upcoming reservations
func (s *TurnOffReminderService) SendTurnOffReminders() {
	s.logger.Info().Msg("SendTurnOffReminders")

	// Find reservations within the next 1 hour that haven't been reminded yet
	reservations, err := s.findReservationsDueForReminder()
	if err != nil {
		s.logger.Err(err).Msg("Failed to fetch reservations for reminders")
		return
	}

	for _, reservation := range reservations {
		// Send reminder and mark as reminded
		if reservation.Product.Post.Status == schema.PostStatusPublished &&
			reservation.Product.Meta.UniWashMachineStatus == schema.UniWashMachineStatusON &&
			reservation.Meta.UniWashLastCommand != "" {
			err := s.processReservationReminder(*reservation)
			if err != nil {
				s.logger.Err(err).Msg("Failed to process reservation reminder")
			}
		}
	}
}

// findReservationsDueForReminder finds reservations due for reminder
func (s *TurnOffReminderService) findReservationsDueForReminder() ([]*schema.Reservation, error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	t := time.Now().In(loc)
	//t := time.Date(2025, 8, 9, 10, 40, 5, 0, loc)

	endTime := t.Add(20 * time.Minute).Truncate(time.Minute)
	startTime := endTime.Add(-1 * time.Hour).Truncate(time.Minute)

	reservations, _, err := s.repo.GetAll(request.Reservations{
		EndTime:   &endTime,
		StartTime: &startTime,
		Status:    schema.ReservationStatusReserved,
	})

	s.logger.Info().Msgf("reservations len %d", len(reservations))

	if err != nil {
		return nil, err
	}

	return reservations, err
}

// processReservationReminder sends SMS and updates reminder status
func (s *TurnOffReminderService) processReservationReminder(reservation schema.Reservation) error {
	if !s.cfg.App.Production {
		return nil
	}

	// Prepare reminder message
	_, err := s.smsService.Send(MessageWay.Message{
		Provider:   5, // با سرشماره 5000
		TemplateID: 16621,
		Method:     "sms",
		Mobile:     fmt.Sprintf("0%d", reservation.User.Mobile),
	})
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}
