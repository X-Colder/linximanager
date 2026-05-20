package model

import "time"

type ProductCategory struct {
	ID         int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID int64    `json:"merchant_id" gorm:"not null;index"`
	Name       string   `json:"name" gorm:"type:varchar(100);not null"`
	ParentID   *int64   `json:"parent_id"`
	SortOrder  int      `json:"sort_order" gorm:"default:0"`
	IconURL    string   `json:"icon_url" gorm:"type:varchar(500)"`
}

func (ProductCategory) TableName() string { return "product_categories" }

type Product struct {
	ID                   int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MerchantID           int64     `json:"merchant_id" gorm:"not null;index:idx_products_merchant"`
	Name                 string    `json:"name" gorm:"type:varchar(200);not null"`
	CategoryID           *int64    `json:"category_id"`
	Description          string    `json:"description" gorm:"type:text"`
	ImageURL             string    `json:"image_url" gorm:"type:varchar(500)"`
	PurchaseUnit         string    `json:"purchase_unit" gorm:"type:varchar(20)"`
	StockUnit            string    `json:"stock_unit" gorm:"type:varchar(20);not null"`
	SaleUnit             string    `json:"sale_unit" gorm:"type:varchar(20);not null"`
	PurchaseToStockRatio float64   `json:"purchase_to_stock_ratio" gorm:"type:decimal(10,4)"`
	StockToSaleRatio     float64   `json:"stock_to_sale_ratio" gorm:"type:decimal(10,4)"`
	CostPrice            float64   `json:"cost_price" gorm:"type:decimal(10,2)"`
	SalePrice            float64   `json:"sale_price" gorm:"type:decimal(10,2);not null"`
	ShelfLifeDays        int       `json:"shelf_life_days"`
	StorageType          string    `json:"storage_type" gorm:"type:varchar(20);default:normal"` // normal/cold/frozen
	LossRate             float64   `json:"loss_rate" gorm:"type:decimal(5,4);default:0"`
	SafetyStock          float64   `json:"safety_stock" gorm:"type:decimal(10,2);default:0"`
	Status               string    `json:"status" gorm:"type:varchar(20);not null;default:active;index:idx_products_merchant"` // active/inactive/deleted
	SortOrder            int       `json:"sort_order" gorm:"default:0"`
	Attributes           JSONMap   `json:"attributes" gorm:"type:jsonb"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// 关联
	Category *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Product) TableName() string { return "products" }

type ProductBOM struct {
	ID         int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID  int64   `json:"product_id" gorm:"not null;index"`  // 成品
	MaterialID int64   `json:"material_id" gorm:"not null;index"` // 原料
	Quantity   float64 `json:"quantity" gorm:"type:decimal(10,4);not null"`
	Unit       string  `json:"unit" gorm:"type:varchar(20);not null"`

	Material *Product `json:"material,omitempty" gorm:"foreignKey:MaterialID"`
}

func (ProductBOM) TableName() string { return "product_bom" }
