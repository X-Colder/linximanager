package repository

import (
	"context"
	"time"

	"github.com/linximanager/backend/internal/model"
	"gorm.io/gorm"
)

type InventoryRepo struct {
	db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
	return &InventoryRepo{db: db}
}

func (r *InventoryRepo) ListByMerchant(ctx context.Context, merchantID int64, offset, limit int) ([]*model.Inventory, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Inventory{}).
		Preload("Product").
		Where("merchant_id = ?", merchantID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*model.Inventory
	err := q.Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *InventoryRepo) FindByProductID(ctx context.Context, merchantID, productID int64) (*model.Inventory, error) {
	var inv model.Inventory
	err := r.db.WithContext(ctx).
		Where("merchant_id = ? AND product_id = ?", merchantID, productID).
		First(&inv).Error
	return &inv, err
}

// Purchase 入库：增加当前库存和可用库存，记录日志（事务内调用）
func (r *InventoryRepo) Purchase(ctx context.Context, tx *gorm.DB, merchantID, productID int64, qty float64, remark string) error {
	var inv model.Inventory
	err := tx.WithContext(ctx).
		Where("merchant_id = ? AND product_id = ?", merchantID, productID).
		First(&inv).Error
	if err == gorm.ErrRecordNotFound {
		inv = model.Inventory{
			MerchantID:     merchantID,
			ProductID:      productID,
			CurrentStock:   qty,
			AvailableStock: qty,
		}
		return tx.WithContext(ctx).Create(&inv).Error
	}
	if err != nil {
		return err
	}
	before := inv.CurrentStock
	if err := tx.WithContext(ctx).Model(&inv).Updates(map[string]any{
		"current_stock":   gorm.Expr("current_stock + ?", qty),
		"available_stock": gorm.Expr("available_stock + ?", qty),
		"updated_at":      time.Now(),
	}).Error; err != nil {
		return err
	}
	log := model.InventoryLog{
		MerchantID: merchantID,
		ProductID:  productID,
		ChangeType: "purchase",
		ChangeQty:  qty,
		BeforeQty:  before,
		AfterQty:   before + qty,
		Remark:     remark,
	}
	return tx.WithContext(ctx).Create(&log).Error
}

// Stocktake 盘点：设置为指定值，记录差异日志（事务内调用）
func (r *InventoryRepo) Stocktake(ctx context.Context, tx *gorm.DB, merchantID, productID int64, newQty float64, remark string) error {
	var inv model.Inventory
	err := tx.WithContext(ctx).
		Where("merchant_id = ? AND product_id = ?", merchantID, productID).
		First(&inv).Error
	if err != nil {
		return err
	}
	before := inv.CurrentStock
	diff := newQty - before
	if err := tx.WithContext(ctx).Model(&inv).Updates(map[string]any{
		"current_stock":   newQty,
		"available_stock": gorm.Expr("available_stock + ?", diff),
		"updated_at":      time.Now(),
	}).Error; err != nil {
		return err
	}
	log := model.InventoryLog{
		MerchantID: merchantID,
		ProductID:  productID,
		ChangeType: "stocktake",
		ChangeQty:  diff,
		BeforeQty:  before,
		AfterQty:   newQty,
		Remark:     remark,
	}
	return tx.WithContext(ctx).Create(&log).Error
}

// DeductAvailable 扣减可用库存（下单锁定）
func (r *InventoryRepo) DeductAvailable(ctx context.Context, tx *gorm.DB, merchantID, productID int64, qty float64) error {
	result := tx.WithContext(ctx).Model(&model.Inventory{}).
		Where("merchant_id = ? AND product_id = ? AND available_stock >= ?", merchantID, productID, qty).
		Updates(map[string]any{
			"available_stock": gorm.Expr("available_stock - ?", qty),
			"locked_stock":    gorm.Expr("locked_stock + ?", qty),
			"updated_at":      time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 库存不足
	}
	return nil
}

// ConsumeLockedStock 核销后：扣减锁定库存和当前库存
func (r *InventoryRepo) ConsumeLockedStock(ctx context.Context, tx *gorm.DB, merchantID, productID int64, qty float64, orderID int64) error {
	if err := tx.WithContext(ctx).Model(&model.Inventory{}).
		Where("merchant_id = ? AND product_id = ?", merchantID, productID).
		Updates(map[string]any{
			"current_stock": gorm.Expr("current_stock - ?", qty),
			"locked_stock":  gorm.Expr("locked_stock - ?", qty),
			"updated_at":    time.Now(),
		}).Error; err != nil {
		return err
	}
	log := model.InventoryLog{
		MerchantID:  merchantID,
		ProductID:   productID,
		ChangeType:  "sale",
		ChangeQty:   -qty,
		ReferenceID: &orderID,
	}
	return tx.WithContext(ctx).Create(&log).Error
}

func (r *InventoryRepo) DB() *gorm.DB {
	return r.db
}
