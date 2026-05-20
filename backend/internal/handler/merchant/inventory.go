package merchant

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/service"
)

type InventoryHandler struct {
	inventorySvc *service.InventoryService
}

func NewInventoryHandler(inventorySvc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventorySvc: inventorySvc}
}

func (h *InventoryHandler) List(c *gin.Context) {
	var q validator.PageParam
	if !validator.BindQuery(c, &q) {
		return
	}
	q.SetDefault()
	mid := middleware.GetMID(c)
	list, total, err := h.inventorySvc.List(c.Request.Context(), mid, q.Offset(), q.PageSize)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OKPage(c, total, q.Page, q.PageSize, list)
}

type purchaseReq struct {
	ProductID int64   `json:"product_id" binding:"required"`
	Qty       float64 `json:"qty" binding:"required,gt=0"`
	Remark    string  `json:"remark"`
}

func (h *InventoryHandler) Purchase(c *gin.Context) {
	var req purchaseReq
	if !validator.BindJSON(c, &req) {
		return
	}
	mid := middleware.GetMID(c)
	if err := h.inventorySvc.Purchase(c.Request.Context(), mid, req.ProductID, req.Qty, req.Remark); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, nil)
}

type stocktakeReq struct {
	ProductID int64   `json:"product_id" binding:"required"`
	NewQty    float64 `json:"new_qty" binding:"min=0"`
	Remark    string  `json:"remark"`
}

func (h *InventoryHandler) Stocktake(c *gin.Context) {
	var req stocktakeReq
	if !validator.BindJSON(c, &req) {
		return
	}
	mid := middleware.GetMID(c)
	if err := h.inventorySvc.Stocktake(c.Request.Context(), mid, req.ProductID, req.NewQty, req.Remark); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, nil)
}
