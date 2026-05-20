package model

import "time"

type Inventory struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID     int64     `json:"merchant_id" gorm:"not null;uniqueIndex:uk_merchant_product"`
	ProductID      int64     `json:"product_id" gorm:"not null;uniqueIndex:uk_merchant_product"`
	CurrentStock   float64   `json:"current_stock" gorm:"type:decimal(10,2);not null;default:0"`
	AvailableStock float64   `json:"available_stock" gorm:"type:decimal(10,2);not null;default:0"`
	LockedStock    float64   `json:"locked_stock" gorm:"type:decimal(10,2);not null;default:0"`
	UpdatedAt      time.Time `json:"updated_at"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (Inventory) TableName() string { return "inventory" }

type InventoryBatch struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	InventoryID    int64     `json:"inventory_id" gorm:"not null;index"`
	BatchNo        string    `json:"batch_no" gorm:"type:varchar(50)"`
	Quantity       float64   `json:"quantity" gorm:"type:decimal(10,2);not null"`
	ProductionDate *time.Time `json:"production_date"`
	ExpiryDate     time.Time `json:"expiry_date" gorm:"not null;index:idx_batches_expiry"`
	Status         string    `json:"status" gorm:"type:varchar(20);not null;default:active;index:idx_batches_expiry"` // active/expired/consumed/damaged
	CreatedAt      time.Time `json:"created_at"`
}

func (InventoryBatch) TableName() string { return "inventory_batches" }

type InventoryLog struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID  int64     `json:"merchant_id" gorm:"not null;index"`
	ProductID   int64     `json:"product_id" gorm:"not null;index"`
	ChangeType  string    `json:"change_type" gorm:"type:varchar(30);not null"` // purchase/sale/loss/adjustment/stocktake
	ChangeQty   float64   `json:"change_qty" gorm:"type:decimal(10,2);not null"`
	BeforeQty   float64   `json:"before_qty" gorm:"type:decimal(10,2);not null"`
	AfterQty    float64   `json:"after_qty" gorm:"type:decimal(10,2);not null"`
	BatchID     *int64    `json:"batch_id"`
	ReferenceID *int64    `json:"reference_id"`
	Remark      string    `json:"remark" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

func (InventoryLog) TableName() string { return "inventory_logs" }
