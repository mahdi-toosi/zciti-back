package response

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/paginator"
	"time"
)

type Transaction struct {
	ID                 uint64
	User               response.User
	Amount             float64
	Status             schema.TransactionStatus
	OrderID            uint64
	UpdatedAt          time.Time
	Description        string
	OrderPaymentMethod schema.OrderPaymentMethod
}

type Transactions struct {
	Data        []*Transaction `json:",omitempty"`
	TotalAmount uint64
	Meta        paginator.Pagination `json:",omitempty"`
}

func FromDomain(item *schema.Transaction) (res *Transaction) {
	if item == nil {
		return nil
	}

	res = &Transaction{
		ID:                 item.ID,
		Amount:             item.Amount,
		Status:             item.Status,
		UpdatedAt:          item.UpdatedAt,
		Description:        item.Description,
		OrderPaymentMethod: item.OrderPaymentMethod,
		User:               response.User{ID: item.User.ID, FullName: item.User.FullName()},
	}

	if item.OrderID != nil {
		res.OrderID = *item.OrderID
	}

	return res
}
