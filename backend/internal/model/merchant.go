package model

import "time"

type Merchant struct {
	ID              int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int64           `json:"user_id" gorm:"not null;index"`
	ShopName        string          `json:"shop_name" gorm:"type:varchar(200);not null"`
	Industry        string          `json:"industry" gorm:"type:varchar(50);not null;index"` // catering/retail/fresh/bakery
	LogoURL         string          `json:"logo_url" gorm:"type:varchar(500)"`
	Address         string          `json:"address" gorm:"type:text"`
	Latitude        float64         `json:"latitude" gorm:"type:decimal(10,7)"`
	Longitude       float64         `json:"longitude" gorm:"type:decimal(10,7)"`
	ContactPhone    string          `json:"contact_phone" gorm:"type:varchar(20)"`
	BusinessHours   string          `json:"business_hours" gorm:"type:varchar(100)"`
	Announcement    string          `json:"announcement" gorm:"type:text"`
	LicenseURL      string          `json:"license_url" gorm:"type:varchar(500)"`
	IDCardURL       string          `json:"id_card_url" gorm:"type:varchar(500)"`
	ShopPhotos      JSONStringSlice `json:"shop_photos" gorm:"type:jsonb"`
	AuditStatus     string          `json:"audit_status" gorm:"type:varchar(20);not null;default:pending"` // pending/approved/rejected
	AuditRemark     string          `json:"audit_remark" gorm:"type:text"`
	VersionPlan     string          `json:"version_plan" gorm:"type:varchar(20);not null;default:basic"` // basic/pro/chain
	PlanExpireAt    *time.Time      `json:"plan_expire_at"`
	MiniappAppid    string          `json:"miniapp_appid" gorm:"type:varchar(100)"`
	MiniappStatus   string          `json:"miniapp_status" gorm:"type:varchar(20);default:pending"` // pending/uploading/auditing/published/rejected
	Status          string          `json:"status" gorm:"type:varchar(20);not null;default:active;index"` // active/frozen
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Merchant) TableName() string { return "merchants" }
