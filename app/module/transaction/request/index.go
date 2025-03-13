package request

import (
	"go-fiber-starter/utils/paginator"
)

//type Transaction struct {
//	ID                   uint64
//	Amount               float64                   `example:"1" validate:"required,number,min=1"`
//	Status               schema.TransactionStatus  `example:"pending" validate:"required"`
//	OrderPaymentMethod   schema.OrderPaymentMethod `example:"online" validate:"required"`
//	GatewayTransactionID string                    `example:"abc" validate:"required"`
//	CreatedAt            time.Time
//}

type Transactions struct {
	WalletID   uint64
	Pagination *paginator.Pagination
}

//
// func (req *Transaction) ToDomain() *schema.Transaction {
//	return &schema.Transaction{
//		ID:                   req.ID,
//		Amount:               req.Amount,
//		Status:               req.Status,
//		OrderPaymentMethod:   req.OrderPaymentMethod,
//		GatewayTransactionID: &req.GatewayTransactionID,
//	}
//}
