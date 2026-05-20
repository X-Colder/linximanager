package repository

import (
	"context"

	"github.com/linximanager/backend/internal/model"
	"gorm.io/gorm"
)

type MerchantRepo struct {
	db *gorm.DB
}

func NewMerchantRepo(db *gorm.DB) *MerchantRepo {
	return &MerchantRepo{db: db}
}

type ListMerchantsFilter struct {
	AuditStatus string
	Status      string
	Industry    string
	Keyword     string
	Offset      int
	Limit       int
}

func (r *MerchantRepo) List(ctx context.Context, f ListMerchantsFilter) ([]*model.Merchant, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Merchant{}).Preload("User")
	if f.AuditStatus != "" {
		q = q.Where("audit_status = ?", f.AuditStatus)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Industry != "" {
		q = q.Where("industry = ?", f.Industry)
	}
	if f.Keyword != "" {
		q = q.Where("shop_name ILIKE ?", "%"+f.Keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var merchants []*model.Merchant
	err := q.Offset(f.Offset).Limit(f.Limit).Order("created_at DESC").Find(&merchants).Error
	return merchants, total, err
}

func (r *MerchantRepo) FindByID(ctx context.Context, id int64) (*model.Merchant, error) {
	var m model.Merchant
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&m).Error
	return &m, err
}

func (r *MerchantRepo) FindByUserID(ctx context.Context, userID int64) (*model.Merchant, error) {
	var m model.Merchant
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error
	return &m, err
}

func (r *MerchantRepo) Create(ctx context.Context, m *model.Merchant) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MerchantRepo) Update(ctx context.Context, m *model.Merchant) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MerchantRepo) UpdateFields(ctx context.Context, id int64, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&model.Merchant{}).Where("id = ?", id).Updates(fields).Error
}

func (r *MerchantRepo) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Merchant{}).Count(&count).Error
	return count, err
}

func (r *MerchantRepo) CountByAuditStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Merchant{}).Where("audit_status = ?", status).Count(&count).Error
	return count, err
}
