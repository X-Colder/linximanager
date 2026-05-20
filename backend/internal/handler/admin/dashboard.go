package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/service"
)

type DashboardHandler struct {
	merchantSvc *service.MerchantService
}

func NewDashboardHandler(merchantSvc *service.MerchantService) *DashboardHandler {
	return &DashboardHandler{merchantSvc: merchantSvc}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	stats, err := h.merchantSvc.PlatformStats(c.Request.Context())
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, stats)
}
