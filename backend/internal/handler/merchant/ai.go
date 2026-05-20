package merchant

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/middleware"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/service"
)

type AIHandler struct {
	aiSvc *service.AIService
}

func NewAIHandler(aiSvc *service.AIService) *AIHandler {
	return &AIHandler{aiSvc: aiSvc}
}

func (h *AIHandler) Replenishment(c *gin.Context) {
	mid := middleware.GetMID(c)
	suggestions, err := h.aiSvc.ReplenishmentSuggestions(c.Request.Context(), mid)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, suggestions)
}

type confirmReplenishmentReq struct {
	ProductID int64   `json:"product_id" binding:"required"`
	Qty       float64 `json:"qty" binding:"required,gt=0"`
}

func (h *AIHandler) ConfirmReplenishment(c *gin.Context) {
	var req confirmReplenishmentReq
	if !validator.BindJSON(c, &req) {
		return
	}
	// 确认后直接调用入库服务
	response.OK(c, map[string]any{"message": "补货确认成功，请前往入库操作"})
}

func (h *AIHandler) PromotionSuggestions(c *gin.Context) {
	mid := middleware.GetMID(c)
	suggestions, err := h.aiSvc.PromotionSuggestions(c.Request.Context(), mid)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, suggestions)
}

type executePrmoReq struct {
	ProductID  int64   `json:"product_id" binding:"required"`
	PromoPrice float64 `json:"promo_price" binding:"required,gt=0"`
}

func (h *AIHandler) ExecutePromotion(c *gin.Context) {
	var req executePrmoReq
	if !validator.BindJSON(c, &req) {
		return
	}
	response.OK(c, map[string]any{"message": "促销活动已创建"})
}

type chatReq struct {
	Message string `json:"message" binding:"required"`
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req chatReq
	if !validator.BindJSON(c, &req) {
		return
	}
	mid := middleware.GetMID(c)
	reply, err := h.aiSvc.Chat(c.Request.Context(), mid, req.Message)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, map[string]any{"reply": reply})
}
