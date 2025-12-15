package test

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/utils/paginator"
	"testing"

	"github.com/rs/zerolog"
)

// =============================================================================
// SendTurnOffReminders Tests
// =============================================================================

func TestTurnOffReminder_SendReminders_Success(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservation := NewReservationBuilder().
		WithID(1).
		WithMobile(9123456789).
		WithPostStatus(schema.PostStatusPublished).
		WithMachineStatus(schema.UniWashMachineStatusON).
		WithLastCommand(schema.UniWashCommandON).
		Build()

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return []*schema.Reservation{reservation}, paginator.Pagination{}, nil
	}

	service.SendTurnOffReminders()

	assertGetAllCalled(t, mockRepo)
	assertMarkTurnOffCalled(t, mockRepo, 1)
}

func TestTurnOffReminder_SendReminders_RepoError(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		Build()

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return nil, paginator.Pagination{}, errors.New("database error")
	}

	service.SendTurnOffReminders()

	assertGetAllCalled(t, mockRepo)
	assertCallCount(t, mockRepo.MarkTurnOffCallCount(), 0, "MarkTurnOffReminderSent")
}

func TestTurnOffReminder_SendReminders_NoReservations(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		Build()

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return []*schema.Reservation{}, paginator.Pagination{}, nil
	}

	service.SendTurnOffReminders()

	assertGetAllCalled(t, mockRepo)
	assertCallCount(t, mockRepo.MarkTurnOffCallCount(), 0, "MarkTurnOffReminderSent")
}

// =============================================================================
// Skip Conditions - Table Driven Tests
// =============================================================================

func TestTurnOffReminder_SendReminders_SkipConditions(t *testing.T) {
	tests := []struct {
		name          string
		postStatus    schema.PostStatus
		machineStatus schema.UniWashMachineStatus
		lastCommand   schema.UniWashCommand
		shouldProcess bool
	}{
		{
			name:          "skips unpublished post",
			postStatus:    schema.PostStatusDraft,
			machineStatus: schema.UniWashMachineStatusON,
			lastCommand:   schema.UniWashCommandON,
			shouldProcess: false,
		},
		{
			name:          "skips machine OFF",
			postStatus:    schema.PostStatusPublished,
			machineStatus: schema.UniWashMachineStatusOFF,
			lastCommand:   schema.UniWashCommandON,
			shouldProcess: false,
		},
		{
			name:          "skips empty last command",
			postStatus:    schema.PostStatusPublished,
			machineStatus: schema.UniWashMachineStatusON,
			lastCommand:   "",
			shouldProcess: false,
		},
		{
			name:          "processes valid reservation",
			postStatus:    schema.PostStatusPublished,
			machineStatus: schema.UniWashMachineStatusON,
			lastCommand:   schema.UniWashCommandON,
			shouldProcess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockCronRepository()
			logger := zerolog.Nop()
			cfg := createTestConfig()

			service := NewTurnOffServiceBuilder().
				WithRepo(mockRepo).
				WithSmsService(cfg, logger).
				Build()

			reservation := NewReservationBuilder().
				WithID(1).
				WithPostStatus(tt.postStatus).
				WithMachineStatus(tt.machineStatus).
				WithLastCommand(tt.lastCommand).
				Build()

			mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
				return []*schema.Reservation{reservation}, paginator.Pagination{}, nil
			}

			service.SendTurnOffReminders()

			if tt.shouldProcess {
				assertMarkTurnOffCalled(t, mockRepo, 1)
			} else {
				assertMarkTurnOffNotCalled(t, mockRepo, 1)
			}
		})
	}
}

// =============================================================================
// Multiple Reservations Tests
// =============================================================================

func TestTurnOffReminder_SendReminders_MultipleReservations(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).WithLastCommand(schema.UniWashCommandON).Build(),
		NewReservationBuilder().WithID(2).WithLastCommand(schema.UniWashCommandON).Build(),
		NewReservationBuilder().WithID(3).WithLastCommand("").Build(), // Should skip - no command
	}

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return reservations, paginator.Pagination{}, nil
	}

	service.SendTurnOffReminders()

	assertMarkTurnOffCalled(t, mockRepo, 1)
	assertMarkTurnOffCalled(t, mockRepo, 2)
	assertMarkTurnOffNotCalled(t, mockRepo, 3)
	assertCallCount(t, mockRepo.MarkTurnOffCallCount(), 2, "MarkTurnOffReminderSent")
}

func TestTurnOffReminder_SendReminders_MixedConditions(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).WithLastCommand(schema.UniWashCommandON).Build(),                                                   // Valid
		NewReservationBuilder().WithID(2).WithPostStatus(schema.PostStatusDraft).WithLastCommand(schema.UniWashCommandON).Build(),            // Invalid - Draft
		NewReservationBuilder().WithID(3).WithMachineStatus(schema.UniWashMachineStatusOFF).WithLastCommand(schema.UniWashCommandON).Build(), // Invalid - OFF
		NewReservationBuilder().WithID(4).WithLastCommand("").Build(),                                                                        // Invalid - no command
		NewReservationBuilder().WithID(5).WithLastCommand(schema.UniWashCommandON).Build(),                                                   // Valid
	}

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return reservations, paginator.Pagination{}, nil
	}

	service.SendTurnOffReminders()

	assertMarkTurnOffCalled(t, mockRepo, 1)
	assertMarkTurnOffNotCalled(t, mockRepo, 2)
	assertMarkTurnOffNotCalled(t, mockRepo, 3)
	assertMarkTurnOffNotCalled(t, mockRepo, 4)
	assertMarkTurnOffCalled(t, mockRepo, 5)
	assertCallCount(t, mockRepo.MarkTurnOffCallCount(), 2, "MarkTurnOffReminderSent")
}

// =============================================================================
// FindReservationsDueForReminder Tests
// =============================================================================

func TestTurnOffReminder_FindReservations(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		Build()

	expectedReservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).WithLastCommand(schema.UniWashCommandON).Build(),
	}

	var capturedReq request.Reservations
	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		capturedReq = req
		return expectedReservations, paginator.Pagination{}, nil
	}

	reservations, err := service.FindReservationsDueForReminder()

	assertNoError(t, err)

	if len(reservations) != 1 {
		t.Errorf("expected 1 reservation, got %d", len(reservations))
	}

	if reservations[0].ID != 1 {
		t.Errorf("expected reservation ID 1, got %d", reservations[0].ID)
	}

	// Verify request parameters
	if capturedReq.TurnOffReminderSent == nil {
		t.Error("expected TurnOffReminderSent to be set")
	} else if *capturedReq.TurnOffReminderSent != false {
		t.Error("expected TurnOffReminderSent to be false")
	}

	if capturedReq.Status != schema.ReservationStatusReserved {
		t.Errorf("expected status to be %s, got %s", schema.ReservationStatusReserved, capturedReq.Status)
	}
}

func TestTurnOffReminder_FindReservations_TimeWindow(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		Build()

	var capturedReq request.Reservations
	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		capturedReq = req
		return []*schema.Reservation{}, paginator.Pagination{}, nil
	}

	_, _ = service.FindReservationsDueForReminder()

	assertTimeWindowIsOneHour(t, capturedReq.StartTime, capturedReq.EndTime)
}

// =============================================================================
// ProcessReservationReminder Tests
// =============================================================================

func TestTurnOffReminder_ProcessReminder_Success(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservation := NewReservationBuilder().
		WithID(1).
		WithLastCommand(schema.UniWashCommandON).
		BuildValue()

	err := service.ProcessReservationReminder(reservation)

	assertNoError(t, err)
	assertMarkTurnOffCalled(t, mockRepo, 1)
}

func TestTurnOffReminder_ProcessReminder_MarkError(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	mockRepo.MarkTurnOffReminderSentFn = func(id uint64) error {
		return errors.New("db error")
	}

	reservation := NewReservationBuilder().
		WithID(1).
		WithLastCommand(schema.UniWashCommandON).
		BuildValue()

	err := service.ProcessReservationReminder(reservation)

	assertError(t, err)
	assertErrorContains(t, err, "failed to mark reminder as sent")
}

// =============================================================================
// Error Resilience Tests
// =============================================================================

func TestTurnOffReminder_SendReminders_ContinuesOnIndividualError(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOffServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).WithLastCommand(schema.UniWashCommandON).Build(),
		NewReservationBuilder().WithID(2).WithLastCommand(schema.UniWashCommandON).Build(),
		NewReservationBuilder().WithID(3).WithLastCommand(schema.UniWashCommandON).Build(),
	}

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return reservations, paginator.Pagination{}, nil
	}

	// Make the second call fail
	mockRepo.MarkTurnOffReminderSentFn = func(id uint64) error {
		if id == 2 {
			return errors.New("simulated error")
		}
		return nil
	}

	service.SendTurnOffReminders()

	// All 3 should be attempted despite error on ID 2
	assertCallCount(t, mockRepo.MarkTurnOffCallCount(), 3, "MarkTurnOffReminderSent")
}
