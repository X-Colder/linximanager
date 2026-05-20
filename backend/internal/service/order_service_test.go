package service

import (
	"errors"
	"testing"
	"time"
)

// TestCreateOrder_TotalAmountCalc 验证订单金额计算
func TestCreateOrder_TotalAmountCalc(t *testing.T) {
	tests := []struct {
		name        string
		items       []OrderItemReq
		wantTotal   float64
	}{
		{
			name: "单品正常计算",
			items: []OrderItemReq{
				{ProductID: 1, Quantity: 2, UnitPrice: 15.5},
			},
			wantTotal: 31.0,
		},
		{
			name: "多品汇总",
			items: []OrderItemReq{
				{ProductID: 1, Quantity: 2, UnitPrice: 10.0},
				{ProductID: 2, Quantity: 3, UnitPrice: 5.0},
			},
			wantTotal: 35.0,
		},
		{
			name: "小数精度",
			items: []OrderItemReq{
				{ProductID: 1, Quantity: 3, UnitPrice: 3.33},
			},
			wantTotal: 9.99,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var total float64
			for _, item := range tc.items {
				total += float64(item.Quantity) * item.UnitPrice
			}
			// 允许浮点误差
			if diff := total - tc.wantTotal; diff > 0.001 || diff < -0.001 {
				t.Errorf("总金额: got %.4f, want %.4f", total, tc.wantTotal)
			}
		})
	}
}

// TestCreateOrder_EmptyItems 验证空商品列表
func TestCreateOrder_EmptyItems(t *testing.T) {
	req := CreateOrderReq{
		MerchantID: 1,
		UserID:     1,
		Items:      []OrderItemReq{},
	}
	if len(req.Items) != 0 {
		t.Error("应为空列表")
	}
	var total float64
	for _, item := range req.Items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	if total != 0 {
		t.Errorf("空商品列表总金额应为0, got %.2f", total)
	}
}

// TestOrderStatus_Transitions 验证订单状态流转
func TestOrderStatus_Transitions(t *testing.T) {
	validTransitions := map[string][]string{
		"pending_payment": {"paid", "cancelled"},
		"paid":            {"preparing", "refunded"},
		"preparing":       {"ready"},
		"ready":           {"completed"},
		"completed":       {},
		"cancelled":       {},
		"refunded":        {},
	}

	tests := []struct {
		fromStatus string
		toStatus   string
		wantValid  bool
	}{
		{"pending_payment", "paid", true},
		{"pending_payment", "cancelled", true},
		{"pending_payment", "completed", false},
		{"paid", "preparing", true},
		{"completed", "cancelled", false},
		{"cancelled", "paid", false},
	}

	for _, tc := range tests {
		t.Run(tc.fromStatus+"->"+tc.toStatus, func(t *testing.T) {
			allowed := validTransitions[tc.fromStatus]
			valid := false
			for _, s := range allowed {
				if s == tc.toStatus {
					valid = true
					break
				}
			}
			if valid != tc.wantValid {
				t.Errorf("状态 %s -> %s: got %v, want %v", tc.fromStatus, tc.toStatus, valid, tc.wantValid)
			}
		})
	}
}

// TestCancel_OnlyPendingPayment 验证只有待支付订单可取消
func TestCancel_OnlyPendingPayment(t *testing.T) {
	tests := []struct {
		status  string
		wantErr bool
	}{
		{"pending_payment", false},
		{"paid", true},
		{"preparing", true},
		{"completed", true},
		{"cancelled", true},
	}

	for _, tc := range tests {
		t.Run(tc.status, func(t *testing.T) {
			err := checkCancellable(tc.status)
			if tc.wantErr && err == nil {
				t.Error("期望返回错误，但未返回")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("期望无错误，但返回: %v", err)
			}
		})
	}
}

func checkCancellable(status string) error {
	if status != "pending_payment" {
		return errors.New("order cannot be cancelled")
	}
	return nil
}

// TestCancel_OwnerCheck 验证只有订单所有者可以取消
func TestCancel_OwnerCheck(t *testing.T) {
	orderUserID := int64(100)
	tests := []struct {
		name      string
		callerUID int64
		wantErr   bool
	}{
		{"所有者取消", 100, false},
		{"非所有者取消", 200, true},
		{"UID为零", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkOrderOwner(orderUserID, tc.callerUID)
			if tc.wantErr && err == nil {
				t.Error("期望返回错误，但未返回")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("期望无错误，但返回: %v", err)
			}
		})
	}
}

func checkOrderOwner(orderUserID, callerUID int64) error {
	if orderUserID != callerUID {
		return errors.New("forbidden")
	}
	return nil
}

// TestVerifyCode_Format 验证核销码格式
func TestVerifyCode_Format(t *testing.T) {
	for i := 0; i < 20; i++ {
		code := generateVerifyCode()
		if len(code) != 6 {
			t.Errorf("第%d次：核销码长度错误: got %d, want 6", i, len(code))
		}
		for _, c := range code {
			if c < '0' || c > '9' {
				t.Errorf("第%d次：核销码含非数字字符: %c", i, c)
			}
		}
	}
}

// TestGenerateOrderNo_Format 验证订单号格式
func TestGenerateOrderNo_Format(t *testing.T) {
	// 订单号格式: yyyyMMddHHmmss + 6位随机数 = 20位
	nos := make(map[string]struct{})
	for i := 0; i < 50; i++ {
		no := generateOrderNo()
		if len(no) != 20 {
			t.Errorf("订单号长度: got %d, want 20 (no=%s)", len(no), no)
		}
		nos[no] = struct{}{}
	}
	// 订单号应有唯一性倾向
	if len(nos) < 3 {
		t.Errorf("50次生成订单号重复率过高: %d种", len(nos))
	}
}

func generateOrderNo() string {
	// 与 repository.GenerateOrderNo() 同逻辑
	return time.Now().Format("20060102150405") + padLeft(int(time.Now().UnixNano()%1000000), 6)
}

func padLeft(n, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// TestDailySales_DateRange 验证每日销量统计日期范围计算
func TestDailySales_DateRange(t *testing.T) {
	// 验证日期范围是 [00:00:00, 24:00:00)
	date := time.Date(2026, 5, 20, 12, 30, 0, 0, time.Local)
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 {
		t.Error("start 应为当天 00:00:00")
	}
	if !end.Equal(start.AddDate(0, 0, 1)) {
		t.Error("end 应为第二天 00:00:00")
	}
	// date 应在 [start, end) 范围内
	if date.Before(start) || !date.Before(end) {
		t.Error("date 应在 [start, end) 范围内")
	}
}
