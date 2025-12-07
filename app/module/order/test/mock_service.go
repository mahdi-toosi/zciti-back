package test

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/order/repository"
	"go-fiber-starter/app/module/order/request"
	"go-fiber-starter/app/module/order/response"
	"go-fiber-starter/utils/paginator"

	"gorm.io/gorm"
)

// MockOrderService is a mock implementation of the order service for testing
type MockOrderService struct {
	repo repository.IRepository
	db   *gorm.DB
}

func (m *MockOrderService) Index(req request.Orders) (orders []*response.Order, paging paginator.Pagination, err error) {
	results, paging, err := m.repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		orders = append(orders, response.FromDomain(result))
	}

	return
}

func (m *MockOrderService) Show(userID uint64, id uint64) (order *response.Order, err error) {
	result, err := m.repo.GetOne(userID, id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result), nil
}

func (m *MockOrderService) Store(req request.Order) (orderID uint64, paymentURL string, err error) {
	// Simplified store for testing - creates order directly without payment processing
	req.Status = schema.OrderStatusPending
	totalAmt := 0.0

	order := req.ToDomain(&totalAmt, nil)
	orderID, err = m.repo.Create(order, nil)
	if err != nil {
		return 0, "", err
	}

	return orderID, "", nil
}

func (m *MockOrderService) Status(userID uint64, orderID uint64, refNum string) (status string, err error) {
	order, err := m.repo.GetOne(userID, orderID)
	if err != nil {
		return "FAILED", err
	}

	order.Status = schema.OrderStatusCompleted
	if err := m.repo.Update(orderID, order); err != nil {
		return "FAILED", err
	}

	return "OK", nil
}

func (m *MockOrderService) Update(id uint64, req request.Order) (err error) {
	return m.repo.Update(id, req.ToDomain(nil, nil))
}

func (m *MockOrderService) Destroy(id uint64) error {
	return m.repo.Delete(id)
}

