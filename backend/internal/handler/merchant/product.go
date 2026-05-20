package merchant

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/model"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/repository"
	"github.com/linximanager/backend/internal/service"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productSvc *service.ProductService
}

func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

type listProductsQuery struct {
	validator.PageParam
	CategoryID *int64 `form:"category_id"`
	Status     string `form:"status"`
	Keyword    string `form:"keyword"`
}

func (h *ProductHandler) List(c *gin.Context) {
	var q listProductsQuery
	if !validator.BindQuery(c, &q) {
		return
	}
	q.SetDefault()
	mid := middleware.GetMID(c)
	products, total, err := h.productSvc.List(c.Request.Context(), repository.ListProductsFilter{
		MerchantID: mid,
		CategoryID: q.CategoryID,
		Status:     q.Status,
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

type createProductReq struct {
	Name                 string         `json:"name" binding:"required,max=200"`
	CategoryID           *int64         `json:"category_id"`
	Description          string         `json:"description"`
	ImageURL             string         `json:"image_url"`
	PurchaseUnit         string         `json:"purchase_unit"`
	StockUnit            string         `json:"stock_unit" binding:"required"`
	SaleUnit             string         `json:"sale_unit" binding:"required"`
	PurchaseToStockRatio float64        `json:"purchase_to_stock_ratio"`
	StockToSaleRatio     float64        `json:"stock_to_sale_ratio"`
	CostPrice            float64        `json:"cost_price"`
	SalePrice            float64        `json:"sale_price" binding:"required,gt=0"`
	ShelfLifeDays        int            `json:"shelf_life_days"`
	StorageType          string         `json:"storage_type"`
	LossRate             float64        `json:"loss_rate"`
	SafetyStock          float64        `json:"safety_stock"`
	Attributes           model.JSONMap  `json:"attributes"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductReq
	if !validator.BindJSON(c, &req) {
		return
	}
	mid := middleware.GetMID(c)
	p := &model.Product{
		MerchantID:           mid,
		Name:                 req.Name,
		CategoryID:           req.CategoryID,
		Description:          req.Description,
		ImageURL:             req.ImageURL,
		PurchaseUnit:         req.PurchaseUnit,
		StockUnit:            req.StockUnit,
		SaleUnit:             req.SaleUnit,
		PurchaseToStockRatio: req.PurchaseToStockRatio,
		StockToSaleRatio:     req.StockToSaleRatio,
		CostPrice:            req.CostPrice,
		SalePrice:            req.SalePrice,
		ShelfLifeDays:        req.ShelfLifeDays,
		StorageType:          req.StorageType,
		LossRate:             req.LossRate,
		SafetyStock:          req.SafetyStock,
		Attributes:           req.Attributes,
		Status:               "active",
	}
	if err := h.productSvc.Create(c.Request.Context(), p); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, p)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseProductID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	mid := middleware.GetMID(c)
	p, err := h.productSvc.GetByID(c.Request.Context(), id, mid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, errcode.ErrProductNotFound)
			return
		}
		response.Fail(c, errcode.ErrInternal)
		return
	}
	var req createProductReq
	if !validator.BindJSON(c, &req) {
		return
	}
	p.Name = req.Name
	p.CategoryID = req.CategoryID
	p.Description = req.Description
	p.ImageURL = req.ImageURL
	p.SalePrice = req.SalePrice
	p.CostPrice = req.CostPrice
	p.SafetyStock = req.SafetyStock
	p.StockUnit = req.StockUnit
	p.SaleUnit = req.SaleUnit
	p.Attributes = req.Attributes
	if err := h.productSvc.Update(c.Request.Context(), p); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseProductID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	mid := middleware.GetMID(c)
	if err := h.productSvc.Delete(c.Request.Context(), id, mid); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, nil)
}

type bomItem struct {
	MaterialID int64   `json:"material_id" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	Unit       string  `json:"unit" binding:"required"`
}

type setBOMReq struct {
	Items []bomItem `json:"items"`
}

func (h *ProductHandler) GetBOM(c *gin.Context) {
	id, err := parseProductID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	bom, err := h.productSvc.GetBOM(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, bom)
}

func (h *ProductHandler) SetBOM(c *gin.Context) {
	id, err := parseProductID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	var req setBOMReq
	if !validator.BindJSON(c, &req) {
		return
	}
	var bom []*model.ProductBOM
	for _, item := range req.Items {
		bom = append(bom, &model.ProductBOM{
			ProductID:  id,
			MaterialID: item.MaterialID,
			Quantity:   item.Quantity,
			Unit:       item.Unit,
		})
	}
	if err := h.productSvc.ReplaceBOM(c.Request.Context(), id, bom); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, nil)
}

func parseProductID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}
