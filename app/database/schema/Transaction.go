package schema

type Transaction struct {
	ID                   uint64             `gorm:"primaryKey"`                             // The unique identifier for the transaction.
	Amount               float64            `gorm:"not null"`                               // The amount of the transaction in the currency specified by the Currency field.
	Status               TransactionStatus  `gorm:"varchar(20); default:pending; not null"` // The current status of the transaction.
	Description          string             `gorm:"varchar(255); not null"`                 //
	OrderPaymentMethod   OrderPaymentMethod `gorm:"default:online; not null"`               // The payment method used for the transaction.
	GatewayTransactionID *string            `gorm:"varchar(255)"`                           // The unique identifier for the transaction provided by the payment gateway.
	WalletID             uint64             `gorm:"index"`                                  // Associated wallet (nullable for non-wallet transactions).
	Wallet               Wallet             `gorm:"foreignKey:WalletID"`                    // Relation to the wallet.
	OrderID              *uint64            ``                                              // The ID of the associated order.
	Order                *Order             `gorm:"foreignKey:OrderID"`                     // The associated order object.
	UserID               uint64             `gorm:"not null"`
	User                 User               `gorm:"foreignKey:UserID"`
	Base
}

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"   // The transaction has been initiated, but has not yet been processed by the payment gateway.
	TransactionStatusSuccess   TransactionStatus = "success"   // The transaction has been successfully processed by the payment gateway.
	TransactionStatusFailed    TransactionStatus = "failed"    // The transaction has failed to process by the payment gateway.
	TransactionStatusRefunded  TransactionStatus = "refunded"  // The transaction has been refunded to the customer.
	TransactionStatusCancelled TransactionStatus = "cancelled" // The transaction has been cancelled by the customer or the administrator.
)
