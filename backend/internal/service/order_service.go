package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/linximanager/backend/internal/model"
	"github.com/linximanager/backend/internal/repository"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo     *repository.OrderRepo
	inventoryRepo *repository.InventoryRepo
}

func NewOrderService(orderRepo *repository.OrderRepo, inventoryRepo *repository.InventoryRepo) *OrderService {
	return &OrderService{orderRepo: orderRepo, inventoryRepo: inventoryRepo}
}

type CreateOrderReq struct {
	MerchantID int64
	UserID     int64
	Items      []OrderItemReq
	CouponID   *int64
	PickupTime *time.Time
}

type OrderItemReq struct {
	ProductID int64
	Name      string
	Quantity  int
	UnitPrice float64
	Spec      string
}

func (s *OrderService) Create(ctx context.Context, req CreateOrderReq) (*model.Order, error) {
	db := s.orderRepo.DB()
	var order *model.Order

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var items []model.OrderItem

		for _, item := range req.Items {
			// 扣减可用库存
			if err := s.inventoryRepo.DeductAvailable(ctx, tx, req.MerchantID, item.ProductID, float64(item.Quantity)); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("inventory insufficient")
				}
				return err
			}
			totalPrice := float64(item.Quantity) * item.UnitPrice
			totalAmount += totalPrice
			items = append(items, model.OrderItem{
				ProductID:   item.ProductID,
				ProductName: item.Name,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				TotalPrice:  totalPrice,
				Spec:        item.Spec,
			})
		}

		order = &model.Order{
			OrderNo:        repository.GenerateOrderNo(),
			MerchantID:     req.MerchantID,
			UserID:         req.UserID,
			Status:         "pending_payment",
			TotalAmount:    totalAmount,
			DiscountAmount: 0,
			PayAmount:      totalAmount,
			CouponID:       req.CouponID,
			PickupTime:     req.PickupTime,
			VerifyCode:     generateVerifyCode(),
			Items:          items,
		}
		return s.orderRepo.Create(ctx, tx, order)
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) List(ctx context.Context, f repository.ListOrdersFilter) ([]*model.Order, int64, error) {
	return s.orderRepo.List(ctx, f)
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}

// Verify 核销订单
func (s *OrderService) Verify(ctx context.Context, merchantID int64, code string) (*model.Order, error) {
	db := s.orderRepo.DB()
	var order *model.Order

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		o, err := s.orderRepo.FindByVerifyCode(ctx, merchantID, code)
		if err != nil {
			return err
		}
		now := time.Now()
		if err := s.orderRepo.UpdateStatus(ctx, tx, o.ID, "completed", map[string]any{
			"verified_at": now,
		}); err != nil {
			return err
		}
		// 扣减实际库存
		for _, item := range o.Items {
			if err := s.inventoryRepo.ConsumeLockedStock(ctx, tx, merchantID, item.ProductID, float64(item.Quantity), o.ID); err != nil {
				return err
			}
		}
		o.Status = "completed"
		o.VerifiedAt = &now
		order = o
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

// MarkPaid 模拟支付回调
func (s *OrderService) MarkPaid(ctx context.Context, orderID int64, wxTransactionID string) error {
	db := s.orderRepo.DB()
	now := time.Now()
	return s.orderRepo.UpdateStatus(ctx, db, orderID, "paid", map[string]any{
		"paid_at":           now,
		"wx_transaction_id": wxTransactionID,
	})
}

// Cancel 取消订单
func (s *OrderService) Cancel(ctx context.Context, orderID, userID int64, reason string) error {
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if o.UserID != userID {
		return errors.New("forbidden")
	}
	if o.Status != "pending_payment" {
		return errors.New("order cannot be cancelled")
	}
	db := s.orderRepo.DB()
	now := time.Now()
	return s.orderRepo.UpdateStatus(ctx, db, orderID, "cancelled", map[string]any{
		"cancelled_at":  now,
		"cancel_reason": reason,
	})
}

func (s *OrderService) DailySales(ctx context.Context, merchantID int64, date time.Time) (float64, int64, error) {
	return s.orderRepo.DailySales(ctx, merchantID, date)
}

func generateVerifyCode() string {
	const chars = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
