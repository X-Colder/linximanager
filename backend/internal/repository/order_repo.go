package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/linximanager/backend/internal/model"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

type ListOrdersFilter struct {
	MerchantID int64
	UserID     int64
	Status     string
	Offset     int
	Limit      int
}

func (r *OrderRepo) List(ctx context.Context, f ListOrdersFilter) ([]*model.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Order{}).Preload("Items")
	if f.MerchantID != 0 {
		q = q.Where("merchant_id = ?", f.MerchantID)
	}
	if f.UserID != 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []*model.Order
	err := q.Offset(f.Offset).Limit(f.Limit).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepo) FindByID(ctx context.Context, id int64) (*model.Order, error) {
	var o model.Order
	err := r.db.WithContext(ctx).Preload("Items").Preload("User").Where("id = ?", id).First(&o).Error
	return &o, err
}

func (r *OrderRepo) FindByVerifyCode(ctx context.Context, merchantID int64, code string) (*model.Order, error) {
	var o model.Order
	err := r.db.WithContext(ctx).Preload("Items").
		Where("merchant_id = ? AND verify_code = ? AND status = 'paid'", merchantID, code).
		First(&o).Error
	return &o, err
}

func (r *OrderRepo) Create(ctx context.Context, tx *gorm.DB, o *model.Order) error {
	return tx.WithContext(ctx).Create(o).Error
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string, extra map[string]any) error {
	fields := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	for k, v := range extra {
		fields[k] = v
	}
	return tx.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Updates(fields).Error
}

// GenerateOrderNo 生成唯一订单号：yyyyMMddHHmmss + 6位随机数
func GenerateOrderNo() string {
	return fmt.Sprintf("%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
}

func (r *OrderRepo) DailySales(ctx context.Context, merchantID int64, date time.Time) (float64, int64, error) {
	type result struct {
		TotalAmount float64
		Count       int64
	}
	var res result
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select("COALESCE(SUM(pay_amount), 0) as total_amount, COUNT(*) as count").
		Where("merchant_id = ? AND status IN ('paid','preparing','ready','completed') AND created_at >= ? AND created_at < ?",
			merchantID, start, end).
		Scan(&res).Error
	return res.TotalAmount, res.Count, err
}

func (r *OrderRepo) DB() *gorm.DB {
	return r.db
}
