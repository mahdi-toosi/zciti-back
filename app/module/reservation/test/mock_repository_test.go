package test

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/reservation/request"
	"go-fiber-starter/app/module/reservation/response"
	"go-fiber-starter/utils/paginator"
	"sync"

	oirequest "go-fiber-starter/app/module/orderItem/request"
)

// MockCronRepository is a mock implementation of repository.IRepository for cron tests
type MockCronRepository struct {
	mu sync.Mutex

	// Function hooks for customizing behavior
	GetAllFunc                func(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error)
	MarkTurnOnReminderSentFn  func(id uint64) error
	MarkTurnOffReminderSentFn func(id uint64) error

	// Call tracking
	getAllCalls                []request.Reservations
	markTurnOnReminderSentIDs  []uint64
	markTurnOffReminderSentIDs []uint64
}

// NewMockCronRepository creates a new mock repository with default behavior
func NewMockCronRepository() *MockCronRepository {
	return &MockCronRepository{
		getAllCalls:                make([]request.Reservations, 0),
		markTurnOnReminderSentIDs:  make([]uint64, 0),
		markTurnOffReminderSentIDs: make([]uint64, 0),
	}
}

// GetAll implements repository.IRepository
func (_m *MockCronRepository) GetAll(req request.Reservations) ([]*schema.Reservation, paginator.Pagination, error) {
	_m.mu.Lock()
	_m.getAllCalls = append(_m.getAllCalls, req)
	_m.mu.Unlock()

	if _m.GetAllFunc != nil {
		return _m.GetAllFunc(req)
	}
	return nil, paginator.Pagination{}, nil
}

// GetOne implements repository.IRepository
func (_m *MockCronRepository) GetOne(businessID uint64, id uint64) (*response.Reservation, error) {
	return nil, nil
}

// Create implements repository.IRepository
func (_m *MockCronRepository) Create(reservation *schema.Reservation) error {
	return nil
}

// Update implements repository.IRepository
func (_m *MockCronRepository) Update(id uint64, reservation *schema.Reservation) error {
	return nil
}

// Delete implements repository.IRepository
func (_m *MockCronRepository) Delete(id uint64) error {
	return nil
}

// IsReservable implements repository.IRepository
func (_m *MockCronRepository) IsReservable(req oirequest.OrderItem, businessID uint64) error {
	return nil
}

// MarkTurnOnReminderSent implements repository.IRepository
func (_m *MockCronRepository) MarkTurnOnReminderSent(id uint64) error {
	_m.mu.Lock()
	_m.markTurnOnReminderSentIDs = append(_m.markTurnOnReminderSentIDs, id)
	_m.mu.Unlock()

	if _m.MarkTurnOnReminderSentFn != nil {
		return _m.MarkTurnOnReminderSentFn(id)
	}
	return nil
}

// MarkTurnOffReminderSent implements repository.IRepository
func (_m *MockCronRepository) MarkTurnOffReminderSent(id uint64) error {
	_m.mu.Lock()
	_m.markTurnOffReminderSentIDs = append(_m.markTurnOffReminderSentIDs, id)
	_m.mu.Unlock()

	if _m.MarkTurnOffReminderSentFn != nil {
		return _m.MarkTurnOffReminderSentFn(id)
	}
	return nil
}

// ===== Assertion Helpers =====

// WasGetAllCalled returns true if GetAll was called at least once
func (_m *MockCronRepository) WasGetAllCalled() bool {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	return len(_m.getAllCalls) > 0
}

// GetAllCallCount returns the number of times GetAll was called
func (_m *MockCronRepository) GetAllCallCount() int {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	return len(_m.getAllCalls)
}

// GetLastGetAllRequest returns the last request passed to GetAll
func (_m *MockCronRepository) GetLastGetAllRequest() *request.Reservations {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	if len(_m.getAllCalls) == 0 {
		return nil
	}
	return &_m.getAllCalls[len(_m.getAllCalls)-1]
}

// WasMarkTurnOnCalled returns true if MarkTurnOnReminderSent was called with the given ID
func (_m *MockCronRepository) WasMarkTurnOnCalled(id uint64) bool {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	for _, v := range _m.markTurnOnReminderSentIDs {
		if v == id {
			return true
		}
	}
	return false
}

// WasMarkTurnOffCalled returns true if MarkTurnOffReminderSent was called with the given ID
func (_m *MockCronRepository) WasMarkTurnOffCalled(id uint64) bool {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	for _, v := range _m.markTurnOffReminderSentIDs {
		if v == id {
			return true
		}
	}
	return false
}

// MarkTurnOnCallCount returns the number of times MarkTurnOnReminderSent was called
func (_m *MockCronRepository) MarkTurnOnCallCount() int {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	return len(_m.markTurnOnReminderSentIDs)
}

// MarkTurnOffCallCount returns the number of times MarkTurnOffReminderSent was called
func (_m *MockCronRepository) MarkTurnOffCallCount() int {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	return len(_m.markTurnOffReminderSentIDs)
}

// GetMarkTurnOnIDs returns all IDs passed to MarkTurnOnReminderSent
func (_m *MockCronRepository) GetMarkTurnOnIDs() []uint64 {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	result := make([]uint64, len(_m.markTurnOnReminderSentIDs))
	copy(result, _m.markTurnOnReminderSentIDs)
	return result
}

// GetMarkTurnOffIDs returns all IDs passed to MarkTurnOffReminderSent
func (_m *MockCronRepository) GetMarkTurnOffIDs() []uint64 {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	result := make([]uint64, len(_m.markTurnOffReminderSentIDs))
	copy(result, _m.markTurnOffReminderSentIDs)
	return result
}

// Reset clears all call tracking data
func (_m *MockCronRepository) Reset() {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	_m.getAllCalls = make([]request.Reservations, 0)
	_m.markTurnOnReminderSentIDs = make([]uint64, 0)
	_m.markTurnOffReminderSentIDs = make([]uint64, 0)
}
