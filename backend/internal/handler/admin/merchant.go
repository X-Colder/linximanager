package admin

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

type MerchantHandler struct {
	merchantSvc *service.MerchantService
}

func NewMerchantHandler(merchantSvc *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{merchantSvc: merchantSvc}
}

type listMerchantsQuery struct {
	validator.PageParam
	AuditStatus string `form:"audit_status"`
	Status      string `form:"status"`
	Industry    string `form:"industry"`
	Keyword     string `form:"keyword"`
}

func (h *MerchantHandler) List(c *gin.Context) {
	var q listMerchantsQuery
	if !validator.BindQuery(c, &q) {
		return
	}
	q.SetDefault()
	merchants, total, err := h.merchantSvc.List(c.Request.Context(), repository.ListMerchantsFilter{
		AuditStatus: q.AuditStatus,
		Status:      q.Status,
		Industry:    q.Industry,
		Keyword:     q.Keyword,
		Offset:      q.Offset(),
		Limit:       q.PageSize,
	})
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OKPage(c, total, q.Page, q.PageSize, merchants)
}

func (h *MerchantHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	m, err := h.merchantSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, errcode.ErrMerchantNotFound)
			return
		}
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, m)
}

type auditReq struct {
	Pass   bool   `json:"pass"`
	Remark string `json:"remark"`
}

func (h *MerchantHandler) Audit(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	var req auditReq
	if !validator.BindJSON(c, &req) {
		return
	}
	if err := h.merchantSvc.Audit(c.Request.Context(), id, req.Pass, req.Remark); err != nil {
		if err.Error() == "merchant not in pending status" {
			response.FailMsg(c, errcode.ErrConflict, err.Error())
			return
		}
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, nil)
}

type freezeReq struct {
	Freeze bool `json:"freeze"`
}

func (h *MerchantHandler) Freeze(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Fail(c, errcode.ErrParamInvalid)
		return
	}
	var req freezeReq
	if !validator.BindJSON(c, &req) {
		return
	}
	if err := h.merchantSvc.Freeze(c.Request.Context(), id, req.Freeze); err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, nil)
}

func parseID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}
