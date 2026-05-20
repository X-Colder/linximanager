package merchant

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/repository"
	"github.com/linximanager/backend/internal/service"
	"strconv"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

type listOrdersQuery struct {
	validator.PageParam
	Status string `form:"status"`
}

func (h *OrderHandler) List(c *gin.Context) {
	var q listOrdersQuery
	if !validator.BindQuery(c, &q) {
		return
	}
	q.SetDefault()
	mid := middleware.GetMID(c)
	orders, total, err := h.orderSvc.List(c.Request.Context(), repository.ListOrdersFilter{
		MerchantID: mid,
		Status:     q.Status,
		Offset:     q.Offset(),
		Limit:      q.PageSize,
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

type verifyReq struct {
	Code string `json:"code" binding:"required"`
}

func (h *OrderHandler) Verify(c *gin.Context) {
	var req verifyReq
	if !validator.BindJSON(c, &req) {
		return
	}
	mid := middleware.GetMID(c)
	order, err := h.orderSvc.Verify(c.Request.Context(), mid, req.Code)
	if err != nil {
		response.FailMsg(c, errcode.ErrVerifyCodeInvalid, err.Error())
		return
	}
	response.OK(c, order)
}
