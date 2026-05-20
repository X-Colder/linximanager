package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/linximanager/backend/internal/model"
	"github.com/linximanager/backend/internal/repository"
)

// AIService 模拟 AI 决策引擎（生产中通过 gRPC 调用 Python 服务）
type AIService struct {
	inventoryRepo *repository.InventoryRepo
	productRepo   *repository.ProductRepo
}

func NewAIService(inventoryRepo *repository.InventoryRepo, productRepo *repository.ProductRepo) *AIService {
	return &AIService{inventoryRepo: inventoryRepo, productRepo: productRepo}
}

// ReplenishmentSuggestions 补货建议（模拟）
func (s *AIService) ReplenishmentSuggestions(ctx context.Context, merchantID int64) ([]*model.ReplenishmentSuggestion, error) {
	invList, _, err := s.inventoryRepo.ListByMerchant(ctx, merchantID, 0, 50)
	if err != nil {
		return nil, err
	}

	var suggestions []*model.ReplenishmentSuggestion
	for _, inv := range invList {
		if inv.Product == nil {
			continue
		}
		// 简单规则：可用库存 < 安全库存 * 1.5 时建议补货
		safetyStock := inv.Product.SafetyStock
		if inv.AvailableStock < safetyStock*1.5 {
			predictedDaily := safetyStock * (0.8 + rand.Float64()*0.4)
			suggestedQty := safetyStock*3 - inv.AvailableStock
			if suggestedQty <= 0 {
				suggestedQty = safetyStock
			}
			suggestions = append(suggestions, &model.ReplenishmentSuggestion{
				MerchantID:           merchantID,
				ProductID:            inv.ProductID,
				CurrentStock:         inv.CurrentStock,
				PredictedDailyDemand: predictedDaily,
				SafetyStock:          safetyStock,
				SuggestedQty:         suggestedQty,
				ExpectedProfit:       suggestedQty * inv.Product.SalePrice * 0.3,
				Status:               "pending",
				CreatedAt:            time.Now(),
				Product:              inv.Product,
			})
		}
	}
	return suggestions, nil
}

// PromotionSuggestions 促销建议（模拟）
func (s *AIService) PromotionSuggestions(ctx context.Context, merchantID int64) ([]map[string]any, error) {
	invList, _, err := s.inventoryRepo.ListByMerchant(ctx, merchantID, 0, 50)
	if err != nil {
		return nil, err
	}

	var suggestions []map[string]any
	for _, inv := range invList {
		if inv.Product == nil {
			continue
		}
		// 模拟：库存过高（超过安全库存5倍）建议促销
		if inv.CurrentStock > inv.Product.SafetyStock*5 && inv.Product.SafetyStock > 0 {
			discountRate := 0.7 + rand.Float64()*0.2
			suggestions = append(suggestions, map[string]any{
				"product_id":       inv.ProductID,
				"product_name":     inv.Product.Name,
				"current_stock":    inv.CurrentStock,
				"original_price":   inv.Product.SalePrice,
				"suggested_price":  inv.Product.SalePrice * discountRate,
				"discount_rate":    discountRate,
				"trigger_reason":   "overstock",
				"predicted_profit": inv.Product.SalePrice * discountRate * inv.CurrentStock * 0.5,
			})
		}
	}
	return suggestions, nil
}

// Chat AI 顾问对话（模拟流式响应）
func (s *AIService) Chat(ctx context.Context, merchantID int64, message string) (string, error) {
	// 生产中此处调用 LLM API 并流式返回
	return "您好，我是灵犀掌柜 AI 顾问。根据您的店铺数据，当前库存健康度良好。如需具体分析，请告诉我您关心的商品或问题。", nil
}
