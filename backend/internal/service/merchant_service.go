package service

import (
	"context"
	"errors"

	"github.com/linximanager/backend/internal/model"
	"github.com/linximanager/backend/internal/repository"
	"gorm.io/gorm"
)

type MerchantService struct {
	merchantRepo *repository.MerchantRepo
	userRepo     *repository.UserRepo
}

func NewMerchantService(merchantRepo *repository.MerchantRepo, userRepo *repository.UserRepo) *MerchantService {
	return &MerchantService{merchantRepo: merchantRepo, userRepo: userRepo}
}

func (s *MerchantService) List(ctx context.Context, f repository.ListMerchantsFilter) ([]*model.Merchant, int64, error) {
	return s.merchantRepo.List(ctx, f)
}

func (s *MerchantService) GetByID(ctx context.Context, id int64) (*model.Merchant, error) {
	return s.merchantRepo.FindByID(ctx, id)
}

func (s *MerchantService) GetByUserID(ctx context.Context, userID int64) (*model.Merchant, error) {
	return s.merchantRepo.FindByUserID(ctx, userID)
}

// Audit 审核商家
func (s *MerchantService) Audit(ctx context.Context, id int64, pass bool, remark string) error {
	m, err := s.merchantRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if m.AuditStatus != "pending" {
		return errors.New("merchant not in pending status")
	}
	status := "approved"
	if !pass {
		status = "rejected"
	}
	fields := map[string]any{
		"audit_status": status,
		"audit_remark": remark,
	}
	if pass {
		// 审核通过时更新 users.role 为 merchant
		_ = s.userRepo.UpdateFields(ctx, m.UserID, map[string]any{"role": "merchant"})
	}
	return s.merchantRepo.UpdateFields(ctx, id, fields)
}

// Freeze 冻结/解冻商家
func (s *MerchantService) Freeze(ctx context.Context, id int64, freeze bool) error {
	status := "active"
	if freeze {
		status = "frozen"
	}
	return s.merchantRepo.UpdateFields(ctx, id, map[string]any{"status": status})
}

// UpdateShop 商家更新自己的店铺信息
func (s *MerchantService) UpdateShop(ctx context.Context, merchantID int64, updates map[string]any) error {
	return s.merchantRepo.UpdateFields(ctx, merchantID, updates)
}

// Register 商家入驻（通常在注册流程中调用）
func (s *MerchantService) Register(ctx context.Context, m *model.Merchant) error {
	return s.merchantRepo.Create(ctx, m)
}

func (s *MerchantService) PlatformStats(ctx context.Context) (map[string]any, error) {
	total, err := s.merchantRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	pending, err := s.merchantRepo.CountByAuditStatus(ctx, "pending")
	if err != nil {
		return nil, err
	}
	approved, err := s.merchantRepo.CountByAuditStatus(ctx, "approved")
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"total_merchants":    total,
		"pending_merchants":  pending,
		"approved_merchants": approved,
	}, nil
}

var _ = gorm.ErrRecordNotFound
