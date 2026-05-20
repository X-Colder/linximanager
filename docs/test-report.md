# 测试报告 — 灵犀掌柜后端单元测试

**日期**: 2026-05-20
**测试人**: qas（自动化测试 Agent）
**状态**: PASS

---

## 一、总体摘要

| 指标 | 数值 |
|------|------|
| 测试覆盖包数 | 5 |
| 测试文件数 | 7 |
| 总测试用例数（含子测试）| 133 |
| 通过 | 133 |
| 失败 | 0 |
| 跳过 | 0 |
| 整体结果 | PASS |

### go test ./... 执行结果

```
ok  github.com/linximanager/backend/internal/handler      0.175s
ok  github.com/linximanager/backend/internal/middleware   0.498s
ok  github.com/linximanager/backend/internal/pkg/jwt      0.495s
ok  github.com/linximanager/backend/internal/pkg/response 0.649s
ok  github.com/linximanager/backend/internal/service      0.727s
```

无编译错误，无测试失败。

---

## 二、分模块测试用例汇总

### 2.1 JWT 工具包（internal/pkg/jwt）

**文件**: `internal/pkg/jwt/jwt_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-JWT-001 | consumer用户正常生成 Access Token | 单元 | PASS |
| TC-JWT-002 | merchant用户含MID生成 Access Token | 单元 | PASS |
| TC-JWT-003 | admin用户生成 Access Token | 单元 | PASS |
| TC-JWT-004 | 错误 secret 解析 Access Token | 单元 | PASS |
| TC-JWT-005 | 过期 Access Token 返回 ErrTokenExpired | 单元 | PASS |
| TC-JWT-006 | 无效 token 字符串返回 ErrTokenInvalid | 单元 | PASS |
| TC-JWT-007 | 空 token 字符串返回 ErrTokenInvalid | 单元 | PASS |
| TC-JWT-008 | 正常生成和解析 Refresh Token | 单元 | PASS |
| TC-JWT-009 | 过期 Refresh Token 返回 ErrTokenExpired | 单元 | PASS |
| TC-JWT-010 | 错误 secret 解析 Refresh Token | 单元 | PASS |
| TC-JWT-011 | 空 secret 生成 Token | 单元 | PASS |
| TC-JWT-012 | ExpiresAt 时间精度验证（±1s）| 单元 | PASS |

**小计**: 12 个用例，12 通过

---

### 2.2 统一响应格式（internal/pkg/response）

**文件**: `internal/pkg/response/response_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-RESP-001 | OK 带数据响应 | 单元 | PASS |
| TC-RESP-002 | OK 空数据 omitempty 省略 data 字段 | 单元 | PASS |
| TC-RESP-003 | OKPage 分页数据格式验证 | 单元 | PASS |
| TC-RESP-004 | 未认证(40100) → HTTP 401 | 单元 | PASS |
| TC-RESP-005 | Token过期(40101) → HTTP 401 | 单元 | PASS |
| TC-RESP-006 | Token无效(40102) → HTTP 401 | 单元 | PASS |
| TC-RESP-007 | 无权限(40300) → HTTP 403 | 单元 | PASS |
| TC-RESP-008 | 商家冻结(40301) → HTTP 403 | 单元 | PASS |
| TC-RESP-009 | 商家待审核(40302) → HTTP 403 | 单元 | PASS |
| TC-RESP-010 | 资源不存在(40400) → HTTP 404 | 单元 | PASS |
| TC-RESP-011 | 用户不存在(40401) → HTTP 404 | 单元 | PASS |
| TC-RESP-012 | 商家不存在(40402) → HTTP 404 | 单元 | PASS |
| TC-RESP-013 | 商品不存在(40403) → HTTP 404 | 单元 | PASS |
| TC-RESP-014 | 订单不存在(40404) → HTTP 404 | 单元 | PASS |
| TC-RESP-015 | 请求过频(42900) → HTTP 429 | 单元 | PASS |
| TC-RESP-016 | 内部错误(50000) → HTTP 500 | 单元 | PASS |
| TC-RESP-017 | 数据库错误(50001) → HTTP 500 | 单元 | PASS |
| TC-RESP-018 | 参数无效(40001) → HTTP 400 | 单元 | PASS |
| TC-RESP-019 | 库存不足(40021) → HTTP 400 | 单元 | PASS |
| TC-RESP-020 | FailMsg 自定义错误消息 | 单元 | PASS |

**小计**: 20 个用例，20 通过

---

### 2.3 认证服务（internal/service/auth_service）

**文件**: `internal/service/auth_service_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-AUTH-001 | buildTokens 生成 access+refresh token | 单元 | PASS |
| TC-AUTH-002 | 手机号不存在时自动注册逻辑验证 | 单元 | PASS |
| TC-AUTH-003 | 冻结账号返回 "account frozen" 错误 | 单元 | PASS |
| TC-AUTH-004 | 正常账号不返回错误 | 单元 | PASS |
| TC-AUTH-005 | 核销码格式为6位数字 | 单元 | PASS |
| TC-AUTH-006 | 核销码100次生成唯一性验证 | 单元 | PASS |
| TC-AUTH-007 | RefreshToken 正确解析 UID | 单元 | PASS |
| TC-AUTH-008 | 过期 RefreshToken 不能刷新 | 单元 | PASS |

**小计**: 8 个用例，8 通过

---

### 2.4 库存服务（internal/service/inventory_service）

**文件**: `internal/service/inventory_service_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-INV-001 | 正常入库100件 | 单元 | PASS |
| TC-INV-002 | 入库0.5件（小数支持）| 单元 | PASS |
| TC-INV-003 | 入库0件（边界值）| 单元 | PASS |
| TC-INV-004 | 入库负数返回错误 | 单元 | PASS |
| TC-INV-005 | 盘点：增加库存差异计算 | 单元 | PASS |
| TC-INV-006 | 盘点：减少库存差异计算 | 单元 | PASS |
| TC-INV-007 | 盘点：相同数量差异为0 | 单元 | PASS |
| TC-INV-008 | 盘点：盘点为零 | 单元 | PASS |
| TC-INV-009 | 盘点：初始为零新增 | 单元 | PASS |
| TC-INV-010 | 扣减可用库存：库存充足 | 单元 | PASS |
| TC-INV-011 | 扣减可用库存：恰好等于扣减量 | 单元 | PASS |
| TC-INV-012 | 扣减可用库存：库存不足 | 单元 | PASS |
| TC-INV-013 | 扣减可用库存：库存为零 | 单元 | PASS |
| TC-INV-014 | 扣减可用库存：扣减量为零 | 单元 | PASS |
| TC-INV-015 | 并发扣减不超卖（20并发，库存100，每次扣10）| 并发 | PASS |
| TC-INV-016 | 核销后 CurrentStock 和 LockedStock 正确减少 | 单元 | PASS |
| TC-INV-017 | 部分核销场景验证 | 单元 | PASS |
| TC-INV-018 | 库存变动类型枚举值校验（purchase/sale/stocktake等）| 单元 | PASS |
| TC-INV-019 | 无效 merchantID/productID 为零的边界判断 | 单元 | PASS |

**小计**: 19 个用例，19 通过

---

### 2.5 订单服务（internal/service/order_service）

**文件**: `internal/service/order_service_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-ORD-001 | 单品订单金额计算 | 单元 | PASS |
| TC-ORD-002 | 多品汇总金额计算 | 单元 | PASS |
| TC-ORD-003 | 小数精度（3×3.33=9.99）| 单元 | PASS |
| TC-ORD-004 | 空商品列表总金额为0 | 单元 | PASS |
| TC-ORD-005 | 订单状态流转：pending→paid（合法）| 单元 | PASS |
| TC-ORD-006 | 订单状态流转：pending→cancelled（合法）| 单元 | PASS |
| TC-ORD-007 | 订单状态流转：pending→completed（非法）| 单元 | PASS |
| TC-ORD-008 | 订单状态流转：paid→preparing（合法）| 单元 | PASS |
| TC-ORD-009 | 订单状态流转：completed→cancelled（非法）| 单元 | PASS |
| TC-ORD-010 | 订单状态流转：cancelled→paid（非法）| 单元 | PASS |
| TC-ORD-011 | 只有 pending_payment 可取消 | 单元 | PASS |
| TC-ORD-012 | paid 状态不可取消 | 单元 | PASS |
| TC-ORD-013 | completed 状态不可取消 | 单元 | PASS |
| TC-ORD-014 | 只有订单所有者可取消 | 单元 | PASS |
| TC-ORD-015 | 非所有者取消返回 forbidden | 单元 | PASS |
| TC-ORD-016 | UID 为零不可取消 | 单元 | PASS |
| TC-ORD-017 | 核销码格式20次验证 | 单元 | PASS |
| TC-ORD-018 | 订单号格式（20位，时间戳+6位数字）| 单元 | PASS |
| TC-ORD-019 | 每日统计日期范围正确（[00:00, 24:00)）| 单元 | PASS |

**小计**: 19 个用例，19 通过

---

### 2.6 认证接口（internal/handler/auth）

**文件**: `internal/handler/auth_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-H-001 | 登录缺少 phone → HTTP 400 | 集成 | PASS |
| TC-H-002 | 登录缺少 code → HTTP 400 | 集成 | PASS |
| TC-H-003 | phone 长度不足11位 → HTTP 400 | 集成 | PASS |
| TC-H-004 | phone 长度超过11位 → HTTP 400 | 集成 | PASS |
| TC-H-005 | 正常手机号登录请求 → HTTP 200 | 集成 | PASS |
| TC-H-006 | 微信登录缺少 code → HTTP 400 | 集成 | PASS |
| TC-H-007 | 微信登录正常请求 → HTTP 200 | 集成 | PASS |
| TC-H-008 | 刷新接口缺少 refresh_token → HTTP 400 | 集成 | PASS |
| TC-H-009 | 登出接口始终返回 200 | 集成 | PASS |
| TC-H-010 | 非法 JSON 请求体 → HTTP 400 | 集成 | PASS |
| TC-H-011 | 登录成功响应包含 access_token 和 refresh_token | 集成 | PASS |

**小计**: 11 个用例，11 通过

---

### 2.7 认证中间件（internal/middleware/auth）

**文件**: `internal/middleware/auth_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-MW-001 | 无 Authorization 头 → 401 | 单元 | PASS |
| TC-MW-002 | 没有 Bearer 前缀 → 401 | 单元 | PASS |
| TC-MW-003 | Bearer 格式缺少空格 → 401 | 单元 | PASS |
| TC-MW-004 | Authorization 为空字符串 → 401 | 单元 | PASS |
| TC-MW-005 | 过期 Token → 401（errcode=40101）| 单元 | PASS |
| TC-MW-006 | 无效 Token → 401（errcode=40102）| 单元 | PASS |
| TC-MW-007 | 有效 Token 通过并注入 uid/role/mid | 单元 | PASS |
| TC-MW-008 | 错误 secret 签名的 Token 被拒绝 | 单元 | PASS |
| TC-MW-009 | admin 角色访问 admin-only 接口 → 200 | 单元 | PASS |
| TC-MW-010 | consumer 角色访问 admin-only 接口 → 403 | 单元 | PASS |
| TC-MW-011 | 多角色允许：merchant 和 staff 均可通过 | 单元 | PASS |
| TC-MW-012 | 多角色：consumer 被拒绝 | 单元 | PASS |
| TC-MW-013 | 多角色：admin 被拒绝（非配置角色）| 单元 | PASS |
| TC-MW-014 | 上下文工具函数未设置时返回零值 | 单元 | PASS |
| TC-MW-015 | 上下文工具函数设置后正确返回 | 单元 | PASS |
| TC-MW-016 | 认证失败后不执行业务 handler | 单元 | PASS |

**小计**: 16 个用例，16 通过

---

### 2.8 限流中间件（internal/middleware/ratelimit）

**文件**: `internal/middleware/ratelimit_test.go`

| 用例编号 | 测试名称 | 类型 | 结果 |
|---------|---------|------|------|
| TC-RL-001 | 低于 burst 的请求全部通过 | 单元 | PASS |
| TC-RL-002 | 超出 burst 后触发 429 限流 | 单元 | PASS |
| TC-RL-003 | 限流响应格式非空 | 单元 | PASS |
| TC-RL-004 | burst=0 时立即返回 429 | 单元 | PASS |
| TC-RL-005 | 高 QPS 配置50并发不误限流 | 并发 | PASS |
| TC-RL-006 | 等待令牌补充后请求可恢复 | 单元 | PASS |
| TC-RL-007 | 30并发下成功数不超过 burst(10) | 并发 | PASS |

**小计**: 7 个用例，7 通过

---

## 三、覆盖率概要

由于 service 和 repository 层依赖真实数据库（PostgreSQL+GORM），当前测试文件采用**纯逻辑单元测试**策略，不连接外部依赖。覆盖情况如下：

| 包 | 策略 | 业务逻辑覆盖率（估算）|
|---|------|-------------------|
| internal/pkg/jwt | 全量逻辑测试 | ~95% |
| internal/pkg/response | 全量逻辑测试 | ~100% |
| internal/pkg/errcode | 通过 response 测试间接覆盖 | ~100% |
| internal/service/auth_service | 核心逻辑测试（无DB）| ~60% |
| internal/service/inventory_service | 核心逻辑测试（无DB）| ~55% |
| internal/service/order_service | 核心逻辑测试（无DB）| ~65% |
| internal/middleware/auth | HTTP 中间件全量测试 | ~95% |
| internal/middleware/ratelimit | 限流全量测试 | ~95% |
| internal/handler/auth | HTTP handler 请求/响应测试 | ~80% |

> 注意：repository 层（CRUD 函数）需要集成测试环境（真实 PostgreSQL）才能完整覆盖，当前未包含在此报告中。

---

## 四、已知问题与观察

### 问题 1：AuthService 依赖具体结构体，不方便 Mock

**位置**: `internal/service/auth_service.go:16-19`

`AuthService` 的字段类型为 `*repository.UserRepo`（具体结构体），而非接口，导致在单元测试中无法直接注入 mock。

**影响**: service 层测试只能验证业务逻辑片段，无法端到端覆盖 LoginByPhone/LoginByWechat/RefreshToken 完整路径。

**建议修复**:
```go
// 定义接口
type UserRepository interface {
    FindByPhone(ctx context.Context, phone string) (*model.User, error)
    FindByOpenid(ctx context.Context, openid string) (*model.User, error)
    FindByID(ctx context.Context, id int64) (*model.User, error)
    Create(ctx context.Context, u *model.User) error
}

// AuthService 改用接口
type AuthService struct {
    userRepo UserRepository
    cfg      *config.Config
}
```

---

### 问题 2：OrderService/InventoryService 同样依赖具体 repo

**位置**: `internal/service/order_service.go:14-17`, `internal/service/inventory_service.go:10-13`

同上，建议为 OrderRepo 和 InventoryRepo 各定义接口，以便单元测试覆盖完整 Create/Verify/Cancel 流程。

---

### 问题 3：generateVerifyCode 使用 math/rand（非密码学安全）

**位置**: `internal/service/order_service.go:163-170`

`math/rand.Intn` 是伪随机，在高并发下可能有碰撞风险（6位数字仅100万种组合）。建议改为：

```go
import "crypto/rand"
import "math/big"

func generateVerifyCode() string {
    const chars = "0123456789"
    b := make([]byte, 6)
    for i := range b {
        n, _ := rand.Int(rand.Reader, big.NewInt(10))
        b[i] = chars[n.Int64()]
    }
    return string(b)
}
```

---

### 问题 4：repository 层缺乏集成测试

当前 `internal/repository/` 下没有测试文件。所有 SQL 操作（Purchase、Stocktake、DeductAvailable、ConsumeLockedStock）未经过实际数据库验证。

**建议**: 使用 testcontainers-go 启动 PostgreSQL 容器进行集成测试：

```go
// 依赖: github.com/testcontainers/testcontainers-go
container, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
    ContainerRequest: testcontainers.ContainerRequest{
        Image:        "postgres:16",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "test_linxi",
        },
    },
    Started: true,
})
```

---

## 五、前端测试现状

前端三个应用（admin、merchant-app、consumer-app）均**未配置 Vitest/Playwright**，当前无可运行的前端测试。

已创建测试计划文档：`/Users/yaojun72/Documents/workspace/llm/linximanager/docs/test-plan.md`，包含：
- 7 个 E2E 测试场景（登录流程、库存操作、下单流程）
- 22 个组件单元测试用例（AlertCard / InventoryShelf / PromoCompare / ProductCard）
- 性能基准指标（首屏 ≤2s，API P50 ≤200ms）

**建议下一步**:
1. 在 `frontend/admin/` 中安装 vitest 和 @vue/test-utils
2. 针对 AlertCard、InventoryShelf、PromoCompare 优先实现单元测试
3. 配置 Playwright 实现管理后台登录流程 E2E 测试

---

## 六、建议优先级

| 优先级 | 建议 | 预期收益 |
|--------|------|---------|
| P0 | 为 UserRepo/OrderRepo/InventoryRepo 定义接口，解锁 service 层 Mock 测试 | service 层覆盖率提升至 ~90% |
| P0 | 添加 repository 集成测试（testcontainers）| 验证真实 SQL 操作正确性 |
| P1 | generateVerifyCode 改为 crypto/rand | 消除碰撞安全风险 |
| P1 | 前端安装 vitest，实现 4 个组件单元测试 | 覆盖核心 UI 组件 |
| P2 | 配置 Playwright E2E，实现登录流程自动化 | 回归测试保障 |
| P2 | 添加 go test -race 检测并发问题 | 验证并发安全性 |
