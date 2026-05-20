package consumer

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"strconv"
)

// CartHandler 基于 Redis Hash 实现购物车
type CartHandler struct {
	rdb *redis.Client
}

func NewCartHandler(rdb *redis.Client) *CartHandler {
	return &CartHandler{rdb: rdb}
}

type CartItem struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	MerchantID  int64   `json:"merchant_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	ImageURL    string  `json:"image_url"`
}

func cartKey(uid int64) string {
	return fmt.Sprintf("cart:%d", uid)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	uid := middleware.GetUID(c)
	ctx := c.Request.Context()
	values, err := h.rdb.HGetAll(ctx, cartKey(uid)).Result()
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	var items []CartItem
	for _, v := range values {
		var item CartItem
		if err := json.Unmarshal([]byte(v), &item); err == nil {
			items = append(items, item)
		}
	}
	response.OK(c, items)
}

type addCartReq struct {
	ProductID   int64   `json:"product_id" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	MerchantID  int64   `json:"merchant_id" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	ImageURL    string  `json:"image_url"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req addCartReq
	if !validator.BindJSON(c, &req) {
		return
	}
	uid := middleware.GetUID(c)
	ctx := c.Request.Context()
	field := fmt.Sprintf("%d", req.ProductID)

	// 若已存在则增加数量
	existing, err := h.rdb.HGet(ctx, cartKey(uid), field).Result()
	var item CartItem
	if err == nil {
		_ = json.Unmarshal([]byte(existing), &item)
		item.Quantity += req.Quantity
	} else {
		item = CartItem{
			ProductID:   req.ProductID,
			ProductName: req.ProductName,
			MerchantID:  req.MerchantID,
			Quantity:    req.Quantity,
			UnitPrice:   req.UnitPrice,
			ImageURL:    req.ImageURL,
		}
	}
	b, _ := json.Marshal(item)
	if err := h.rdb.HSet(ctx, cartKey(uid), field, string(b)).Err(); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, item)
}

type updateCartReq struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	uid := middleware.GetUID(c)
	pid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	var req updateCartReq
	if !validator.BindJSON(c, &req) {
		return
	}
	ctx := c.Request.Context()
	field := fmt.Sprintf("%d", pid)
	if req.Quantity == 0 {
		h.rdb.HDel(ctx, cartKey(uid), field)
		response.OK(c, nil)
		return
	}
	existing, err := h.rdb.HGet(ctx, cartKey(uid), field).Result()
	if err != nil {
		response.Fail(c, errcode.ErrNotFound)
		return
	}
	var item CartItem
	_ = json.Unmarshal([]byte(existing), &item)
	item.Quantity = req.Quantity
	b, _ := json.Marshal(item)
	h.rdb.HSet(ctx, cartKey(uid), field, string(b))
	response.OK(c, item)
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	uid := middleware.GetUID(c)
	pid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	ctx := c.Request.Context()
	h.rdb.HDel(ctx, cartKey(uid), fmt.Sprintf("%d", pid))
	response.OK(c, nil)
}
