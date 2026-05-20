package merchant

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/service"
)

type DashboardHandler struct {
	merchantSvc  *service.MerchantService
	orderSvc     *service.OrderService
	inventorySvc *service.InventoryService
}

func NewDashboardHandler(merchantSvc *service.MerchantService, orderSvc *service.OrderService, inventorySvc *service.InventoryService) *DashboardHandler {
	return &DashboardHandler{merchantSvc: merchantSvc, orderSvc: orderSvc, inventorySvc: inventorySvc}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	mid := middleware.GetMID(c)
	ctx := c.Request.Context()

	todaySales, todayOrders, err := h.orderSvc.DailySales(ctx, mid, time.Now())
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}

	invList, _, err := h.inventorySvc.List(ctx, mid, 0, 100)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}

	var lowStockCount int
	for _, inv := range invList {
		if inv.Product != nil && inv.AvailableStock < inv.Product.SafetyStock {
			lowStockCount++
		}
	}

	response.OK(c, map[string]any{
		"today_sales":     todaySales,
		"today_orders":    todayOrders,
		"low_stock_count": lowStockCount,
		"total_products":  len(invList),
	})
}

func (h *DashboardHandler) Alerts(c *gin.Context) {
	mid := middleware.GetMID(c)
	ctx := c.Request.Context()

	invList, _, err := h.inventorySvc.List(ctx, mid, 0, 200)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}

	var alerts []map[string]any
	for _, inv := range invList {
		if inv.Product == nil {
			continue
		}
		if inv.AvailableStock < inv.Product.SafetyStock {
			alerts = append(alerts, map[string]any{
				"type":            "low_stock",
				"product_id":      inv.ProductID,
				"product_name":    inv.Product.Name,
				"available_stock": inv.AvailableStock,
				"safety_stock":    inv.Product.SafetyStock,
			})
		}
	}
	response.OK(c, alerts)
}

func (h *DashboardHandler) Todos(c *gin.Context) {
	response.OK(c, []map[string]any{})
}
