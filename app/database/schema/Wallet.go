package schema

type Wallet struct {
	ID         uint64   `gorm:"primaryKey"`            // The unique identifier for the transaction.
	Amount     float64  `gorm:""`                      // The amount of the transaction
	UserID     *uint64  `gorm:"index:idx_wallet"`      // The ID of the associated user.
	User       User     `gorm:"foreignKey:UserID"`     // The associated user object.
	BusinessID *uint64  `gorm:"index:idx_wallet"`      // The ID of the associated business.
	Business   Business `gorm:"foreignKey:BusinessID"` // The associated business object.
	Base
}
