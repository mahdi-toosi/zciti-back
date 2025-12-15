package test

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/reservation/cron"
	"go-fiber-starter/internal"
	"go-fiber-starter/utils/config"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// ===== Config Helpers =====

// createTestConfig creates a test config with fake API key
func createTestConfig() *config.Config {
	cfg := &config.Config{}
	cfg.Services.MessageWay.ApiKey = "test-api-key"
	return cfg
}

// ===== Reservation Factory =====

// ReservationBuilder provides a fluent API for building test reservations
type ReservationBuilder struct {
	reservation *schema.Reservation
}

// NewReservationBuilder creates a new builder with sensible defaults
func NewReservationBuilder() *ReservationBuilder {
	return &ReservationBuilder{
		reservation: &schema.Reservation{
			ID:         1,
			UserID:     1,
			ProductID:  1,
			BusinessID: 1,
			Status:     schema.ReservationStatusReserved,
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(time.Hour),
			User: schema.User{
				ID:     1,
				Mobile: 9123456789,
			},
			Product: schema.Product{
				ID: 1,
				Post: schema.Post{
					ID:     1,
					Status: schema.PostStatusPublished,
				},
				Meta: schema.ProductMeta{
					UniWashMachineStatus: schema.UniWashMachineStatusON,
				},
			},
			Meta: schema.ReservationMeta{
				UniWashLastCommand: "",
			},
		},
	}
}

func (_b *ReservationBuilder) WithID(id uint64) *ReservationBuilder {
	_b.reservation.ID = id
	return _b
}

func (_b *ReservationBuilder) WithMobile(mobile uint64) *ReservationBuilder {
	_b.reservation.User.Mobile = mobile
	return _b
}

func (_b *ReservationBuilder) WithPostStatus(status schema.PostStatus) *ReservationBuilder {
	_b.reservation.Product.Post.Status = status
	return _b
}

func (_b *ReservationBuilder) WithMachineStatus(status schema.UniWashMachineStatus) *ReservationBuilder {
	_b.reservation.Product.Meta.UniWashMachineStatus = status
	return _b
}

func (_b *ReservationBuilder) WithLastCommand(cmd schema.UniWashCommand) *ReservationBuilder {
	_b.reservation.Meta.UniWashLastCommand = cmd
	return _b
}

func (_b *ReservationBuilder) Build() *schema.Reservation {
	return _b.reservation
}

func (_b *ReservationBuilder) BuildValue() schema.Reservation {
	return *_b.reservation
}

// ===== Service Factory =====

// TurnOnServiceBuilder provides a fluent API for building TurnOnReminderService
type TurnOnServiceBuilder struct {
	cfg      *config.Config
	logger   zerolog.Logger
	repo     *MockCronRepository
	cronSpec string
}

func NewTurnOnServiceBuilder() *TurnOnServiceBuilder {
	return &TurnOnServiceBuilder{
		cfg:      createTestConfig(),
		logger:   zerolog.Nop(),
		cronSpec: "@every 1m",
	}
}

func (_b *TurnOnServiceBuilder) WithRepo(repo *MockCronRepository) *TurnOnServiceBuilder {
	_b.repo = repo
	return _b
}

func (_b *TurnOnServiceBuilder) WithSmsService(cfg *config.Config, logger zerolog.Logger) *TurnOnServiceBuilder {
	_b.cfg = cfg
	_b.logger = logger
	return _b
}

func (_b *TurnOnServiceBuilder) WithConfig(cfg *config.Config) *TurnOnServiceBuilder {
	_b.cfg = cfg
	return _b
}

func (_b *TurnOnServiceBuilder) Build() *cron.TurnOnReminderService {
	return &cron.TurnOnReminderService{
		Cfg:        _b.cfg,
		Logger:     _b.logger,
		Repo:       _b.repo,
		CronSpec:   _b.cronSpec,
		SmsService: internal.NewMessageWay(_b.cfg, _b.logger),
	}
}

// TurnOffServiceBuilder provides a fluent API for building TurnOffReminderService
type TurnOffServiceBuilder struct {
	cfg      *config.Config
	logger   zerolog.Logger
	repo     *MockCronRepository
	cronSpec string
}

func NewTurnOffServiceBuilder() *TurnOffServiceBuilder {
	return &TurnOffServiceBuilder{
		cfg:      createTestConfig(),
		logger:   zerolog.Nop(),
		cronSpec: "@every 1m",
	}
}

func (_b *TurnOffServiceBuilder) WithRepo(repo *MockCronRepository) *TurnOffServiceBuilder {
	_b.repo = repo
	return _b
}

func (_b *TurnOffServiceBuilder) WithSmsService(cfg *config.Config, logger zerolog.Logger) *TurnOffServiceBuilder {
	_b.cfg = cfg
	_b.logger = logger
	return _b
}

func (_b *TurnOffServiceBuilder) WithConfig(cfg *config.Config) *TurnOffServiceBuilder {
	_b.cfg = cfg
	return _b
}

func (_b *TurnOffServiceBuilder) Build() *cron.TurnOffReminderService {
	return &cron.TurnOffReminderService{
		Cfg:        _b.cfg,
		Logger:     _b.logger,
		Repo:       _b.repo,
		CronSpec:   _b.cronSpec,
		SmsService: internal.NewMessageWay(_b.cfg, _b.logger),
	}
}

// ===== Test Assertions =====

// assertMarkTurnOnCalled asserts that MarkTurnOnReminderSent was called with the given ID
func assertMarkTurnOnCalled(t *testing.T, mock *MockCronRepository, id uint64) {
	t.Helper()
	if !mock.WasMarkTurnOnCalled(id) {
		t.Errorf("expected MarkTurnOnReminderSent to be called with ID %d", id)
	}
}

// assertMarkTurnOnNotCalled asserts that MarkTurnOnReminderSent was NOT called with the given ID
func assertMarkTurnOnNotCalled(t *testing.T, mock *MockCronRepository, id uint64) {
	t.Helper()
	if mock.WasMarkTurnOnCalled(id) {
		t.Errorf("expected MarkTurnOnReminderSent NOT to be called with ID %d", id)
	}
}

// assertMarkTurnOffCalled asserts that MarkTurnOffReminderSent was called with the given ID
func assertMarkTurnOffCalled(t *testing.T, mock *MockCronRepository, id uint64) {
	t.Helper()
	if !mock.WasMarkTurnOffCalled(id) {
		t.Errorf("expected MarkTurnOffReminderSent to be called with ID %d", id)
	}
}

// assertMarkTurnOffNotCalled asserts that MarkTurnOffReminderSent was NOT called with the given ID
func assertMarkTurnOffNotCalled(t *testing.T, mock *MockCronRepository, id uint64) {
	t.Helper()
	if mock.WasMarkTurnOffCalled(id) {
		t.Errorf("expected MarkTurnOffReminderSent NOT to be called with ID %d", id)
	}
}

// assertCallCount asserts the expected call count
func assertCallCount(t *testing.T, got, want int, methodName string) {
	t.Helper()
	if got != want {
		t.Errorf("expected %d calls to %s, got %d", want, methodName, got)
	}
}

// assertNoError asserts that no error occurred
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// assertError asserts that an error occurred
func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// assertErrorContains asserts that the error message contains the expected substring
func assertErrorContains(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing '%s', got nil", substr)
		return
	}
	if len(err.Error()) < len(substr) || err.Error()[:len(substr)] != substr {
		t.Errorf("expected error to start with '%s', got '%s'", substr, err.Error())
	}
}

// assertGetAllCalled asserts that GetAll was called
func assertGetAllCalled(t *testing.T, mock *MockCronRepository) {
	t.Helper()
	if !mock.WasGetAllCalled() {
		t.Error("expected GetAll to be called")
	}
}

// assertTimeWindowIsOneHour asserts that the time window between StartTime and EndTime is 1 hour
func assertTimeWindowIsOneHour(t *testing.T, startTime, endTime *time.Time) {
	t.Helper()
	if startTime == nil {
		t.Fatal("expected StartTime to be set")
	}
	if endTime == nil {
		t.Fatal("expected EndTime to be set")
	}
	duration := endTime.Sub(*startTime)
	if duration != time.Hour {
		t.Errorf("expected time window of 1 hour, got %v", duration)
	}
}
