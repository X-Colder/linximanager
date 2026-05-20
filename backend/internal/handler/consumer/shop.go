package consumer

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/repository"
	"github.com/linximanager/backend/internal/service"
	"gorm.io/gorm"
)

type ShopHandler struct {
	merchantSvc *service.MerchantService
	productSvc  *service.ProductService
}

func NewShopHandler(merchantSvc *service.MerchantService, productSvc *service.ProductService) *ShopHandler {
	return &ShopHandler{merchantSvc: merchantSvc, productSvc: productSvc}
}

func (h *ShopHandler) GetShop(c *gin.Context) {
	mid, err := strconv.ParseInt(c.Param("merchant_id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	m, err := h.merchantSvc.GetByID(c.Request.Context(), mid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, errcode.ErrMerchantNotFound)
			return
		}
		response.Fail(c, errcode.ErrInternal)
		return
	}
	if m.Status != "active" || m.AuditStatus != "approved" {
		response.Fail(c, errcode.ErrMerchantNotFound)
		return
	}
	response.OK(c, m)
}

type listShopProductsQuery struct {
	validator.PageParam
	CategoryID *int64 `form:"category_id"`
	Keyword    string `form:"keyword"`
}

func (h *ShopHandler) ListProducts(c *gin.Context) {
	mid, err := strconv.ParseInt(c.Param("merchant_id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	var q listShopProductsQuery
	if !validator.BindQuery(c, &q) {
		return
	}
	q.SetDefault()
	products, total, err := h.productSvc.List(c.Request.Context(), repository.ListProductsFilter{
		MerchantID: mid,
		CategoryID: q.CategoryID,
		Status:     "active",
		Keyword:    q.Keyword,
		Offset:     q.Offset(),
		Limit:      q.PageSize,
	})
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OKPage(c, total, q.Page, q.PageSize, products)
}

func (h *ShopHandler) GetProduct(c *gin.Context) {
	pid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	p, err := h.productSvc.GetPublicByID(c.Request.Context(), pid)
	if err != nil {
		response.Fail(c, errcode.ErrProductNotFound)
		return
	}
	response.OK(c, p)
}
