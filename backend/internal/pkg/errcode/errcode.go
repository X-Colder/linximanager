package errcode

const (
	// 通用成功
	Success = 0

	// 4xxxx 业务错误
	ErrParamInvalid    = 40001 // 参数校验失败
	ErrParamMissing    = 40002 // 参数缺失
	ErrUnauthorized    = 40100 // 未认证
	ErrTokenExpired    = 40101 // Token已过期
	ErrTokenInvalid    = 40102 // Token无效
	ErrForbidden       = 40300 // 无权限
	ErrNotFound        = 40400 // 资源不存在
	ErrConflict        = 40900 // 资源冲突
	ErrTooManyRequests = 42900 // 请求过于频繁

	// 业务场景错误
	ErrUserNotFound       = 40401 // 用户不存在
	ErrPasswordIncorrect  = 40011 // 密码错误
	ErrSMSCodeInvalid     = 40012 // 验证码错误
	ErrMerchantNotFound   = 40402 // 商家不存在
	ErrMerchantFrozen     = 40301 // 商家已冻结
	ErrMerchantAuditPend  = 40302 // 商家待审核
	ErrProductNotFound    = 40403 // 商品不存在
	ErrStockInsufficient  = 40021 // 库存不足
	ErrOrderNotFound      = 40404 // 订单不存在
	ErrOrderStatusInvalid = 40022 // 订单状态不允许此操作
	ErrVerifyCodeInvalid  = 40023 // 核销码无效

	// 5xxxx 系统错误
	ErrInternal    = 50000 // 内部错误
	ErrDB          = 50001 // 数据库错误
	ErrCache       = 50002 // 缓存错误
	ErrExternalAPI = 50003 // 外部API错误
)

var messages = map[int]string{
	Success:               "success",
	ErrParamInvalid:       "参数校验失败",
	ErrParamMissing:       "参数缺失",
	ErrUnauthorized:       "未认证，请先登录",
	ErrTokenExpired:       "Token已过期，请重新登录",
	ErrTokenInvalid:       "Token无效",
	ErrForbidden:          "无权限访问",
	ErrNotFound:           "资源不存在",
	ErrConflict:           "资源已存在",
	ErrTooManyRequests:    "请求过于频繁，请稍后再试",
	ErrUserNotFound:       "用户不存在",
	ErrPasswordIncorrect:  "密码错误",
	ErrSMSCodeInvalid:     "验证码错误或已过期",
	ErrMerchantNotFound:   "商家不存在",
	ErrMerchantFrozen:     "商家账号已冻结",
	ErrMerchantAuditPend:  "商家审核尚未通过",
	ErrProductNotFound:    "商品不存在",
	ErrStockInsufficient:  "库存不足",
	ErrOrderNotFound:      "订单不存在",
	ErrOrderStatusInvalid: "当前订单状态不允许此操作",
	ErrVerifyCodeInvalid:  "核销码无效或已使用",
	ErrInternal:           "内部服务错误",
	ErrDB:                 "数据库错误",
	ErrCache:              "缓存服务错误",
	ErrExternalAPI:        "外部服务调用失败",
}

func Message(code int) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
