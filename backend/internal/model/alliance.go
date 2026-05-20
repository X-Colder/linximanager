package model

import "time"

type AllianceMember struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID   int64     `json:"merchant_id" gorm:"not null;uniqueIndex"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null;default:active"`
	PromoCredits int       `json:"promo_credits" gorm:"default:0"`
	JoinedAt     time.Time `json:"joined_at" gorm:"not null;default:now()"`
}

func (AllianceMember) TableName() string { return "alliance_members" }

type AllianceExposure struct {
	ID               int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	PromotionID      int64  `json:"promotion_id" gorm:"not null;index;uniqueIndex:uk_exposure"`
	SourceMerchantID int64  `json:"source_merchant_id" gorm:"not null;uniqueIndex:uk_exposure"`
	TargetMerchantID int64  `json:"target_merchant_id" gorm:"not null;index"`
	Impressions      int    `json:"impressions" gorm:"default:0"`
	Clicks           int    `json:"clicks" gorm:"default:0"`
	Conversions      int    `json:"conversions" gorm:"default:0"`
	Date             string `json:"date" gorm:"type:date;not null;uniqueIndex:uk_exposure"`
}

func (AllianceExposure) TableName() string { return "alliance_exposures" }

type ReplenishmentSuggestion struct {
	ID                   int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID           int64     `json:"merchant_id" gorm:"not null;index"`
	ProductID            int64     `json:"product_id" gorm:"not null;index"`
	CurrentStock         float64   `json:"current_stock" gorm:"type:decimal(10,2);not null"`
	PredictedDailyDemand float64   `json:"predicted_daily_demand" gorm:"type:decimal(10,2);not null"`
	SafetyStock          float64   `json:"safety_stock" gorm:"type:decimal(10,2);not null"`
	SuggestedQty         float64   `json:"suggested_qty" gorm:"type:decimal(10,2);not null"`
	ExpectedProfit       float64   `json:"expected_profit" gorm:"type:decimal(10,2)"`
	Status               string    `json:"status" gorm:"type:varchar(20);not null;default:pending"` // pending/accepted/rejected
	CreatedAt            time.Time `json:"created_at"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (ReplenishmentSuggestion) TableName() string { return "replenishment_suggestions" }
