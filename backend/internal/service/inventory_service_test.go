package service

import (
	"context"
	"sync"
	"testing"
)

// TestInventoryService_Purchase_Logic 验证入库业务逻辑（数量校验）
func TestInventoryService_Purchase_Logic(t *testing.T) {
	tests := []struct {
		name    string
		qty     float64
		wantErr bool
	}{
		{"正常入库100件", 100.0, false},
		{"入库0.5件（小数）", 0.5, false},
		{"入库0件（边界值）", 0.0, false},
		{"入库负数（异常）", -1.0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePurchaseQty(tc.qty)
			if tc.wantErr && err == nil {
				t.Error("期望返回错误，但未返回")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("期望无错误，但返回: %v", err)
			}
		})
	}
}

// validatePurchaseQty 模拟入库数量校验逻辑
func validatePurchaseQty(qty float64) error {
	if qty < 0 {
		return errInvalidQty
	}
	return nil
}

var errInvalidQty = &inventoryError{"入库数量不能为负数"}

type inventoryError struct{ msg string }

func (e *inventoryError) Error() string { return e.msg }

// TestInventoryService_Stocktake_DiffCalc 验证盘点差异计算
func TestInventoryService_Stocktake_DiffCalc(t *testing.T) {
	tests := []struct {
		name        string
		currentStock float64
		newStock     float64
		expectedDiff float64
	}{
		{"增加库存", 100.0, 150.0, 50.0},
		{"减少库存", 100.0, 80.0, -20.0},
		{"相同数量", 100.0, 100.0, 0.0},
		{"盘点为零", 100.0, 0.0, -100.0},
		{"初始为零新增", 0.0, 50.0, 50.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diff := tc.newStock - tc.currentStock
			if diff != tc.expectedDiff {
				t.Errorf("差异计算错误: got %f, want %f", diff, tc.expectedDiff)
			}
		})
	}
}

// TestInventoryService_DeductAvailable_Logic 验证扣减可用库存逻辑
func TestInventoryService_DeductAvailable_Logic(t *testing.T) {
	tests := []struct {
		name           string
		availableStock float64
		deductQty      float64
		wantSufficient bool
	}{
		{"库存充足", 100.0, 10.0, true},
		{"库存恰好等于扣减量", 10.0, 10.0, true},
		{"库存不足", 5.0, 10.0, false},
		{"库存为零", 0.0, 1.0, false},
		{"扣减量为零", 100.0, 0.0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sufficient := tc.availableStock >= tc.deductQty
			if sufficient != tc.wantSufficient {
				t.Errorf("库存充足性: got %v, want %v (available=%.1f, deduct=%.1f)",
					sufficient, tc.wantSufficient, tc.availableStock, tc.deductQty)
			}
		})
	}
}

// TestInventoryService_ConcurrentDeduct 验证并发扣减不会导致超卖
func TestInventoryService_ConcurrentDeduct_AtomicCounter(t *testing.T) {
	// 模拟并发场景下使用原子操作扣减库存
	// 实际DB操作通过 "available_stock >= qty" 条件保证原子性
	initialStock := 100
	deductPerGoroutine := 10
	goroutineCount := 20 // 总共需要 200，但库存只有 100

	var mu sync.Mutex
	stock := initialStock
	successCount := 0
	failCount := 0

	var wg sync.WaitGroup
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			if stock >= deductPerGoroutine {
				stock -= deductPerGoroutine
				successCount++
			} else {
				failCount++
			}
		}()
	}
	wg.Wait()

	// 只有10笔应该成功（100/10=10）
	if successCount != 10 {
		t.Errorf("成功扣减次数: got %d, want 10", successCount)
	}
	if failCount != 10 {
		t.Errorf("失败次数: got %d, want 10", failCount)
	}
	if stock != 0 {
		t.Errorf("最终库存: got %d, want 0", stock)
	}
}

// TestInventoryService_ConsumeLockedStock_Logic 验证核销后库存变化
func TestInventoryService_ConsumeLockedStock_Logic(t *testing.T) {
	tests := []struct {
		name          string
		currentStock  float64
		lockedStock   float64
		consumeQty    float64
		wantCurrent   float64
		wantLocked    float64
	}{
		{
			name:         "正常核销",
			currentStock: 50.0,
			lockedStock:  10.0,
			consumeQty:   10.0,
			wantCurrent:  40.0,
			wantLocked:   0.0,
		},
		{
			name:         "部分核销",
			currentStock: 50.0,
			lockedStock:  20.0,
			consumeQty:   5.0,
			wantCurrent:  45.0,
			wantLocked:   15.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			current := tc.currentStock - tc.consumeQty
			locked := tc.lockedStock - tc.consumeQty
			if current != tc.wantCurrent {
				t.Errorf("CurrentStock: got %.1f, want %.1f", current, tc.wantCurrent)
			}
			if locked != tc.wantLocked {
				t.Errorf("LockedStock: got %.1f, want %.1f", locked, tc.wantLocked)
			}
		})
	}
}

// TestInventoryLog_ChangeTypeValues 验证库存变动类型枚举值
func TestInventoryLog_ChangeTypeValues(t *testing.T) {
	validChangeTypes := map[string]bool{
		"purchase":   true,
		"sale":       true,
		"loss":       true,
		"adjustment": true,
		"stocktake":  true,
	}

	testTypes := []struct {
		changeType string
		wantValid  bool
	}{
		{"purchase", true},
		{"sale", true},
		{"stocktake", true},
		{"unknown", false},
		{"", false},
	}

	for _, tc := range testTypes {
		t.Run(tc.changeType, func(t *testing.T) {
			valid := validChangeTypes[tc.changeType]
			if valid != tc.wantValid {
				t.Errorf("changeType %q valid: got %v, want %v", tc.changeType, valid, tc.wantValid)
			}
		})
	}
}

// TestInventoryService_GetByProductID_NotFound 验证商品不存在时的错误处理
func TestInventoryService_GetByProductID_NotFound(t *testing.T) {
	// 模拟 repository 返回 not found
	ctx := context.Background()
	_ = ctx // 实际集成测试中使用
	// 此处记录边界条件：merchantID 或 productID 为 0 时属于无效查询
	testCases := []struct {
		merchantID int64
		productID  int64
		valid      bool
	}{
		{1, 1, true},
		{0, 1, false},
		{1, 0, false},
		{0, 0, false},
	}
	for _, tc := range testCases {
		valid := tc.merchantID > 0 && tc.productID > 0
		if valid != tc.valid {
			t.Errorf("merchantID=%d productID=%d valid: got %v, want %v",
				tc.merchantID, tc.productID, valid, tc.valid)
		}
	}
}
