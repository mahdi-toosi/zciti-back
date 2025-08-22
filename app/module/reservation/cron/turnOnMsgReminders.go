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

type TurnOnReminderService struct {
	cronSpec   string
	cfg        *config.Config
	logger     zerolog.Logger
	smsService *MessageWay.App // Interface for SMS sending
	repo       repository.IRepository
}

func RunTurnOnMsgReminders(
	cfg *config.Config,
	logger zerolog.Logger,
	Repo repository.IRepository,
	smsService *MessageWay.App,
	cronService *internal.CronService,
) *TurnOnReminderService {
	service := &TurnOnReminderService{
		cfg:        cfg,
		repo:       Repo,
		logger:     logger,
		smsService: smsService,
		cronSpec:   "@every 10s",
	}

	err := cronService.AddJob(service.cronSpec, service.SendReservationReminders)
	if err != nil {
		service.logger.Fatal().Err(err).Msg("failed to add RunTurnOnMsgReminders job")
	}

	return service
}

// SendReservationReminders checks and sends reminders for upcoming reservations
func (s *TurnOnReminderService) SendReservationReminders() {
	s.logger.Info().Msg("SendReservationReminders")

	// Find reservations within the next 1 hour that haven't been reminded yet
	reservations, err := s.findReservationsDueForReminder()
	if err != nil {
		s.logger.Err(err).Msg("Failed to fetch reservations for reminders")
		return
	}

	for _, reservation := range reservations {
		// Send reminder and mark as reminded
		if reservation.Product.Post.Status == schema.PostStatusPublished &&
			reservation.Product.Meta.UniWashMachineStatus == schema.UniWashMachineStatusON {
			err := s.processReservationReminder(*reservation)
			if err != nil {
				s.logger.Err(err).Msg("Failed to process reservation reminder")
			}
		}
	}
}

// findReservationsDueForReminder finds reservations due for reminder
func (s *TurnOnReminderService) findReservationsDueForReminder() ([]*schema.Reservation, error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	t := time.Now().In(loc)
	//t := time.Date(2025, 8, 9, 10, 00, 5, 0, loc)

	startTime := t.Add(60 * time.Minute).Truncate(time.Minute)
	endTime := startTime.Add(time.Hour)

	reservations, _, err := s.repo.GetAll(request.Reservations{
		EndTime:   &endTime,
		StartTime: &startTime,
		Status:    schema.ReservationStatusReserved,
	})
	if err != nil {
		return nil, err
	}

	return reservations, err
}

// processReservationReminder sends SMS and updates reminder status
func (s *TurnOnReminderService) processReservationReminder(reservation schema.Reservation) error {
	if !s.cfg.App.Production {
		return nil
	}

	// Prepare reminder message
	_, err := s.smsService.Send(MessageWay.Message{
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
