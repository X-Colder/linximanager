package model

import "time"

type Order struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo         string     `json:"order_no" gorm:"type:varchar(32);uniqueIndex;not null"`
	MerchantID      int64      `json:"merchant_id" gorm:"not null;index:idx_orders_merchant_status"`
	UserID          int64      `json:"user_id" gorm:"not null;index:idx_orders_user"`
	Status          string     `json:"status" gorm:"type:varchar(20);not null;default:pending_payment;index:idx_orders_merchant_status"`
	// pending_payment/paid/preparing/ready/completed/cancelled/refunded
	TotalAmount     float64    `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	DiscountAmount  float64    `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	PayAmount       float64    `json:"pay_amount" gorm:"type:decimal(10,2);not null"`
	CouponID        *int64     `json:"coupon_id"`
	PickupTime      *time.Time `json:"pickup_time"`
	VerifyCode      string     `json:"verify_code,omitempty" gorm:"type:varchar(10);index:idx_orders_verify"`
	VerifiedAt      *time.Time `json:"verified_at"`
	PaidAt          *time.Time `json:"paid_at"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	CancelReason    string     `json:"cancel_reason" gorm:"type:text"`
	WxTransactionID string     `json:"wx_transaction_id,omitempty" gorm:"type:varchar(100)"`
	CreatedAt       time.Time  `json:"created_at" gorm:"index:idx_orders_user"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	User  *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID          int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID     int64   `json:"order_id" gorm:"not null;index"`
	ProductID   int64   `json:"product_id" gorm:"not null"`
	ProductName string  `json:"product_name" gorm:"type:varchar(200);not null"`
	Quantity    int     `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice  float64 `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Spec        string  `json:"spec" gorm:"type:varchar(100)"`
}

func (OrderItem) TableName() string { return "order_items" }
