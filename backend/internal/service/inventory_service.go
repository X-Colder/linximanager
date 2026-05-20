package service

import (
	"context"

	"github.com/linximanager/backend/internal/model"
	"github.com/linximanager/backend/internal/repository"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepo
}

func NewInventoryService(inventoryRepo *repository.InventoryRepo) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) List(ctx context.Context, merchantID int64, offset, limit int) ([]*model.Inventory, int64, error) {
	return s.inventoryRepo.ListByMerchant(ctx, merchantID, offset, limit)
}

// Purchase 入库
func (s *InventoryService) Purchase(ctx context.Context, merchantID, productID int64, qty float64, remark string) error {
	db := s.inventoryRepo.DB()
	tx := db.WithContext(ctx).Begin()
	if err := s.inventoryRepo.Purchase(ctx, tx, merchantID, productID, qty, remark); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Stocktake 盘点
func (s *InventoryService) Stocktake(ctx context.Context, merchantID, productID int64, newQty float64, remark string) error {
	db := s.inventoryRepo.DB()
	tx := db.WithContext(ctx).Begin()
	if err := s.inventoryRepo.Stocktake(ctx, tx, merchantID, productID, newQty, remark); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *InventoryService) GetByProductID(ctx context.Context, merchantID, productID int64) (*model.Inventory, error) {
	return s.inventoryRepo.FindByProductID(ctx, merchantID, productID)
}
