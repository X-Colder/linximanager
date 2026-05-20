package model

import "time"

type Promotion struct {
	ID               int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID       int64      `json:"merchant_id" gorm:"not null;index"`
	ProductID        int64      `json:"product_id" gorm:"not null;index"`
	Type             string     `json:"type" gorm:"type:varchar(20);not null"` // clearance/flash_sale/new_arrival
	TriggerType      string     `json:"trigger_type" gorm:"type:varchar(20)"`  // auto/manual
	TriggerReason    string     `json:"trigger_reason" gorm:"type:varchar(20)"` // overstock/expiring
	OriginalPrice    float64    `json:"original_price" gorm:"type:decimal(10,2);not null"`
	PromoPrice       float64    `json:"promo_price" gorm:"type:decimal(10,2);not null"`
	DiscountRate     float64    `json:"discount_rate" gorm:"type:decimal(4,2)"`
	PredictedSellDays int       `json:"predicted_sell_days"`
	PredictedProfit  float64    `json:"predicted_profit" gorm:"type:decimal(10,2)"`
	Status           string     `json:"status" gorm:"type:varchar(20);not null;default:pending"` // pending/active/completed/cancelled
	StartAt          *time.Time `json:"start_at"`
	EndAt            *time.Time `json:"end_at"`
	ActualSoldQty    int        `json:"actual_sold_qty" gorm:"default:0"`
	ActualProfit     float64    `json:"actual_profit" gorm:"type:decimal(10,2)"`
	AllianceEnabled  bool       `json:"alliance_enabled" gorm:"default:false"`
	CreatedAt        time.Time  `json:"created_at"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (Promotion) TableName() string { return "promotions" }

type Coupon struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID int64     `json:"merchant_id" gorm:"not null;index"`
	Name       string    `json:"name" gorm:"type:varchar(100);not null"`
	Type       string    `json:"type" gorm:"type:varchar(20);not null"` // fixed/percent
	Value      float64   `json:"value" gorm:"type:decimal(10,2);not null"`
	MinAmount  float64   `json:"min_amount" gorm:"type:decimal(10,2);default:0"`
	TotalQty   int       `json:"total_qty" gorm:"not null"`
	UsedQty    int       `json:"used_qty" gorm:"default:0"`
	StartAt    time.Time `json:"start_at" gorm:"not null"`
	EndAt      time.Time `json:"end_at" gorm:"not null"`
	Status     string    `json:"status" gorm:"type:varchar(20);not null;default:active"`
}

func (Coupon) TableName() string { return "coupons" }

type UserCoupon struct {
	ID       int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID   int64      `json:"user_id" gorm:"not null;index"`
	CouponID int64      `json:"coupon_id" gorm:"not null;index"`
	Status   string     `json:"status" gorm:"type:varchar(20);not null;default:unused"` // unused/used/expired
	UsedAt   *time.Time `json:"used_at"`
	OrderID  *int64     `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`

	Coupon *Coupon `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
}

func (UserCoupon) TableName() string { return "user_coupons" }
