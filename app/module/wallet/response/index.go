package response

import (
	"go-fiber-starter/app/database/schema"
)

type Wallet struct {
	ID         uint64
	Amount     float64
	BusinessID *uint64
	UserID     *uint64
}

func FromDomain(item *schema.Wallet) (res *Wallet) {
	if item == nil {
		return nil
	}

	res = &Wallet{
		ID:         item.ID,
		Amount:     item.Amount,
		UserID:     item.UserID,
		BusinessID: item.BusinessID,
	}

	return res
}
