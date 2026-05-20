package consumer

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/repository"
	"github.com/linximanager/backend/internal/service"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

type createOrderReq struct {
	MerchantID int64         `json:"merchant_id" binding:"required"`
	Items      []orderItemReq `json:"items" binding:"required,min=1"`
	CouponID   *int64        `json:"coupon_id"`
	PickupTime *time.Time    `json:"pickup_time"`
}

type orderItemReq struct {
	ProductID int64   `json:"product_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,gt=0"`
	Spec      string  `json:"spec"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderReq
	if !validator.BindJSON(c, &req) {
		return
	}
	uid := middleware.GetUID(c)
	var items []service.OrderItemReq
	for _, item := range req.Items {
		items = append(items, service.OrderItemReq{
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Spec:      item.Spec,
		})
	}
	order, err := h.orderSvc.Create(c.Request.Context(), service.CreateOrderReq{
		MerchantID: req.MerchantID,
		UserID:     uid,
		Items:      items,
		CouponID:   req.CouponID,
		PickupTime: req.PickupTime,
	})
	if err != nil {
		if err.Error() == "inventory insufficient" {
			response.Fail(c, errcode.ErrStockInsufficient)
			return
		}
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, order)
}

type listMyOrdersQuery struct {
	validator.PageParam
	Status string `form:"status"`
}

func (h *OrderHandler) List(c *gin.Context) {
	var q listMyOrdersQuery
	if !validator.BindQuery(c, &q) {
		return
	}
	q.SetDefault()
	uid := middleware.GetUID(c)
	orders, total, err := h.orderSvc.List(c.Request.Context(), repository.ListOrdersFilter{
		UserID: uid,
		Status: q.Status,
		Offset: q.Offset(),
		Limit:  q.PageSize,
	})
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OKPage(c, total, q.Page, q.PageSize, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	order, err := h.orderSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, errcode.ErrOrderNotFound)
		return
	}
	response.OK(c, order)
}

func (h *OrderHandler) Pay(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	// 模拟支付成功回调
	if err := h.orderSvc.MarkPaid(c.Request.Context(), id, "mock_wx_tx_"+strconv.FormatInt(id, 10)); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, map[string]any{"message": "支付成功"})
}

type cancelReq struct {
	Reason string `json:"reason"`
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	var req cancelReq
	_ = c.ShouldBindJSON(&req)
	uid := middleware.GetUID(c)
	if err := h.orderSvc.Cancel(c.Request.Context(), id, uid, req.Reason); err != nil {
		response.FailMsg(c, errcode.ErrOrderStatusInvalid, err.Error())
		return
	}
	response.OK(c, nil)
}
