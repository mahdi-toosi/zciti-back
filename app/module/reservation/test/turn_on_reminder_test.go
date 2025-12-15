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
// SendTurnOnReminders Tests
// =============================================================================

func TestTurnOnReminder_SendReminders_Success(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservation := NewReservationBuilder().
		WithID(1).
		WithMobile(9123456789).
		WithPostStatus(schema.PostStatusPublished).
		WithMachineStatus(schema.UniWashMachineStatusON).
		Build()

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return []*schema.Reservation{reservation}, paginator.Pagination{}, nil
	}

	service.SendTurnOnReminders()

	assertGetAllCalled(t, mockRepo)
	assertMarkTurnOnCalled(t, mockRepo, 1)
}

func TestTurnOnReminder_SendReminders_RepoError(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		Build()

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return nil, paginator.Pagination{}, errors.New("database error")
	}

	service.SendTurnOnReminders()

	assertGetAllCalled(t, mockRepo)
	assertCallCount(t, mockRepo.MarkTurnOnCallCount(), 0, "MarkTurnOnReminderSent")
}

func TestTurnOnReminder_SendReminders_NoReservations(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		Build()

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return []*schema.Reservation{}, paginator.Pagination{}, nil
	}

	service.SendTurnOnReminders()

	assertGetAllCalled(t, mockRepo)
	assertCallCount(t, mockRepo.MarkTurnOnCallCount(), 0, "MarkTurnOnReminderSent")
}

// =============================================================================
// Skip Conditions - Table Driven Tests
// =============================================================================

func TestTurnOnReminder_SendReminders_SkipConditions(t *testing.T) {
	tests := []struct {
		name          string
		postStatus    schema.PostStatus
		machineStatus schema.UniWashMachineStatus
		shouldProcess bool
	}{
		{
			name:          "skips unpublished post",
			postStatus:    schema.PostStatusDraft,
			machineStatus: schema.UniWashMachineStatusON,
			shouldProcess: false,
		},
		{
			name:          "skips machine OFF",
			postStatus:    schema.PostStatusPublished,
			machineStatus: schema.UniWashMachineStatusOFF,
			shouldProcess: false,
		},
		{
			name:          "processes valid reservation",
			postStatus:    schema.PostStatusPublished,
			machineStatus: schema.UniWashMachineStatusON,
			shouldProcess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockCronRepository()
			logger := zerolog.Nop()
			cfg := createTestConfig()

			service := NewTurnOnServiceBuilder().
				WithRepo(mockRepo).
				WithSmsService(cfg, logger).
				Build()

			reservation := NewReservationBuilder().
				WithID(1).
				WithPostStatus(tt.postStatus).
				WithMachineStatus(tt.machineStatus).
				Build()

			mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
				return []*schema.Reservation{reservation}, paginator.Pagination{}, nil
			}

			service.SendTurnOnReminders()

			if tt.shouldProcess {
				assertMarkTurnOnCalled(t, mockRepo, 1)
			} else {
				assertMarkTurnOnNotCalled(t, mockRepo, 1)
			}
		})
	}
}

// =============================================================================
// Multiple Reservations Tests
// =============================================================================

func TestTurnOnReminder_SendReminders_MultipleReservations(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).WithMobile(9123456789).Build(),
		NewReservationBuilder().WithID(2).WithMobile(9987654321).Build(),
		NewReservationBuilder().WithID(3).WithMobile(9111222333).WithPostStatus(schema.PostStatusDraft).Build(), // Should skip
	}

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return reservations, paginator.Pagination{}, nil
	}

	service.SendTurnOnReminders()

	assertMarkTurnOnCalled(t, mockRepo, 1)
	assertMarkTurnOnCalled(t, mockRepo, 2)
	assertMarkTurnOnNotCalled(t, mockRepo, 3)
	assertCallCount(t, mockRepo.MarkTurnOnCallCount(), 2, "MarkTurnOnReminderSent")
}

func TestTurnOnReminder_SendReminders_MixedConditions(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).Build(),                                                   // Valid
		NewReservationBuilder().WithID(2).WithPostStatus(schema.PostStatusDraft).Build(),            // Invalid - Draft
		NewReservationBuilder().WithID(3).WithMachineStatus(schema.UniWashMachineStatusOFF).Build(), // Invalid - OFF
		NewReservationBuilder().WithID(4).Build(),                                                   // Valid
		NewReservationBuilder().WithID(5).WithPostStatus(schema.PostStatusDraft).Build(),            // Invalid - Draft
	}

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return reservations, paginator.Pagination{}, nil
	}

	service.SendTurnOnReminders()

	assertMarkTurnOnCalled(t, mockRepo, 1)
	assertMarkTurnOnNotCalled(t, mockRepo, 2)
	assertMarkTurnOnNotCalled(t, mockRepo, 3)
	assertMarkTurnOnCalled(t, mockRepo, 4)
	assertMarkTurnOnNotCalled(t, mockRepo, 5)
	assertCallCount(t, mockRepo.MarkTurnOnCallCount(), 2, "MarkTurnOnReminderSent")
}

// =============================================================================
// FindReservationsDueForReminder Tests
// =============================================================================

func TestTurnOnReminder_FindReservations(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		Build()

	expectedReservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).Build(),
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
	if capturedReq.TurnOnReminderSent == nil {
		t.Error("expected TurnOnReminderSent to be set")
	} else if *capturedReq.TurnOnReminderSent != false {
		t.Error("expected TurnOnReminderSent to be false")
	}

	if capturedReq.Status != schema.ReservationStatusReserved {
		t.Errorf("expected status to be %s, got %s", schema.ReservationStatusReserved, capturedReq.Status)
	}
}

func TestTurnOnReminder_FindReservations_TimeWindow(t *testing.T) {
	mockRepo := NewMockCronRepository()
	service := NewTurnOnServiceBuilder().
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

func TestTurnOnReminder_ProcessReminder_Success(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservation := NewReservationBuilder().WithID(1).BuildValue()

	err := service.ProcessReservationReminder(reservation)

	assertNoError(t, err)
	assertMarkTurnOnCalled(t, mockRepo, 1)
}

func TestTurnOnReminder_ProcessReminder_MarkError(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	mockRepo.MarkTurnOnReminderSentFn = func(id uint64) error {
		return errors.New("db error")
	}

	reservation := NewReservationBuilder().WithID(1).BuildValue()

	err := service.ProcessReservationReminder(reservation)

	assertError(t, err)
	assertErrorContains(t, err, "failed to mark reminder as sent")
}

// =============================================================================
// Error Resilience Tests
// =============================================================================

func TestTurnOnReminder_SendReminders_ContinuesOnIndividualError(t *testing.T) {
	mockRepo := NewMockCronRepository()
	logger := zerolog.Nop()
	cfg := createTestConfig()

	service := NewTurnOnServiceBuilder().
		WithRepo(mockRepo).
		WithSmsService(cfg, logger).
		Build()

	reservations := []*schema.Reservation{
		NewReservationBuilder().WithID(1).Build(),
		NewReservationBuilder().WithID(2).Build(),
		NewReservationBuilder().WithID(3).Build(),
	}

	mockRepo.GetAllFunc = func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
		return reservations, paginator.Pagination{}, nil
	}

	// Make the second call fail
	mockRepo.MarkTurnOnReminderSentFn = func(id uint64) error {
		if id == 2 {
			return errors.New("simulated error")
		}
		return nil
	}

	service.SendTurnOnReminders()

	// All 3 should be attempted despite error on ID 2
	assertCallCount(t, mockRepo.MarkTurnOnCallCount(), 3, "MarkTurnOnReminderSent")
}
