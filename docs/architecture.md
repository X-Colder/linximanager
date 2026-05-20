# 灵犀掌柜 - 智能库存决策平台 系统架构设计文档

## 1. 技术选型

### 1.1 前端技术栈

| 技术 | 选型 | 理由 |
|:---|:---|:---|
| 框架 | Vue 3 + TypeScript | 生态成熟、类型安全、Composition API适合复杂逻辑 |
| 跨端方案 | uni-app | 一套代码生成微信小程序(商家端+C端)和Web端 |
| 状态管理 | Pinia | Vue 3官方推荐，轻量、TypeScript友好 |
| UI组件 | Element Plus (Web) / uni-ui (小程序) | Web端企业级组件，小程序端原生性能 |
| 图表 | ECharts | 数据可视化丰富、移动端适配好 |
| HTTP | Axios | 拦截器、请求取消、TypeScript支持 |
| 构建 | Vite | 快速HMR、Tree-shaking |

### 1.2 后端技术栈

| 技术 | 选型 | 理由 |
|:---|:---|:---|
| 语言 | Go 1.22+ | 高并发、低延迟、部署简单 |
| Web框架 | Gin | 高性能、中间件丰富、社区活跃 |
| ORM | GORM | Go生态主流，支持Migration |
| 数据库 | PostgreSQL 16 | JSONB支持复杂商品属性、分区表优化大数据查询 |
| 缓存 | Redis 7 | 高性能KV、支持Pub/Sub、Stream |
| 消息队列 | Redis Stream + NATS | 轻量高效，NATS用于服务间通信 |
| 对象存储 | MinIO / 阿里云OSS | 图片、证照、模板资源 |
| 搜索 | Meilisearch | 轻量全文搜索、商品搜索 |
| AI推理 | Python微服务 (FastAPI) | AI生态丰富，通过gRPC与Go通信 |
| 定时任务 | Asynq (基于Redis) | Go原生、支持延迟任务和定时任务 |

### 1.3 基础设施

| 技术 | 选型 | 理由 |
|:---|:---|:---|
| 容器 | Docker | 标准化部署 |
| 编排 | Docker Compose (开发) / K8s (生产) | 灵活扩缩容 |
| CI/CD | GitHub Actions | 免费额度充足、生态集成好 |
| 网关 | Nginx / Traefik | 反向代理、SSL、负载均衡 |
| 监控 | Prometheus + Grafana | 业界标准、告警灵活 |
| 日志 | Loki + Promtail | 轻量、与Grafana集成 |
| 链路追踪 | OpenTelemetry + Jaeger | 微服务调用链分析 |

---

## 2. 系统架构

### 2.1 整体分层架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端层 (Client Layer)                    │
│  ┌──────────┐   ┌──────────────┐   ┌──────────────────┐         │
│  │ Web管理台 │   │ 商家小程序    │   │   C端用户小程序   │         │
│  │ (Vue3)   │   │ (uni-app)    │   │   (uni-app)      │         │
│  └──────────┘   └──────────────┘   └──────────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      网关层 (Gateway Layer)                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  Nginx/Traefik (SSL终止、限流、路由、负载均衡)            │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    服务层 (Service Layer)                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │ 用户服务 │ │ 商家服务 │ │ 商品服务 │ │ 订单服务 │ │ 支付服务 │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │ AI决策  │ │ 联盟服务 │ │ 促销服务 │ │ 通知服务 │ │ 文件服务 │  │
│  │ 引擎    │ │         │ │         │ │         │ │         │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     数据层 (Data Layer)                           │
│  ┌──────────┐  ┌───────┐  ┌──────────┐  ┌──────────────────┐   │
│  │PostgreSQL│  │ Redis │  │  MinIO   │  │   Meilisearch    │   │
│  │(主从复制)│  │(集群) │  │(对象存储)│  │   (全文搜索)     │   │
│  └──────────┘  └───────┘  └──────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 微服务划分

采用**模块化单体**架构起步，按领域边界清晰划分模块，未来可按需拆分为独立微服务：

| 服务模块 | 职责 | 对外端口 |
|:---|:---|:---|
| gateway | API网关、路由、限流、认证 | 8080 |
| user-service | 用户注册登录、角色权限、商家入驻审核 | 内部 |
| merchant-service | 商家档案、店铺管理、版本管理 | 内部 |
| product-service | 商品CRUD、BOM管理、批次管理、库存管理 | 内部 |
| order-service | 订单创建、状态流转、核销 | 内部 |
| payment-service | 微信支付、退款、对账 | 内部 |
| ai-engine | 库存预测、促销建议、补货决策、AI顾问 | gRPC 50051 |
| promotion-service | 促销方案管理、优惠券、限时活动 | 内部 |
| alliance-service | 流量联盟、曝光分发、推广券 | 内部 |
| notification-service | 短信、模板消息、站内信 | 内部 |
| file-service | 图片上传、证照存储、模板包管理 | 内部 |
| admin-service | 平台管理后台专用接口 | 内部 |

---

## 3. 数据库Schema设计

### 3.1 核心表结构

```sql
-- ==================== 用户与认证 ====================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'consumer', -- admin/merchant/staff/consumer
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active/frozen/pending
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    openid VARCHAR(100) UNIQUE,
    unionid VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==================== 商家 ====================

CREATE TABLE merchants (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    shop_name VARCHAR(200) NOT NULL,
    industry VARCHAR(50) NOT NULL, -- catering/retail/fresh/bakery
    logo_url VARCHAR(500),
    address TEXT,
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    contact_phone VARCHAR(20),
    business_hours VARCHAR(100),
    announcement TEXT,
    license_url VARCHAR(500),
    id_card_url VARCHAR(500),
    shop_photos JSONB, -- ["url1","url2"]
    audit_status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending/approved/rejected
    audit_remark TEXT,
    version_plan VARCHAR(20) NOT NULL DEFAULT 'basic', -- basic/pro/chain
    plan_expire_at TIMESTAMPTZ,
    miniapp_appid VARCHAR(100),
    miniapp_status VARCHAR(20) DEFAULT 'pending', -- pending/uploading/auditing/published/rejected
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active/frozen
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_merchants_industry ON merchants(industry);
CREATE INDEX idx_merchants_status ON merchants(status);
CREATE INDEX idx_merchants_location ON merchants USING GIST (
    ST_MakePoint(longitude, latitude)
);

-- ==================== 商品 ====================

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    name VARCHAR(200) NOT NULL,
    category_id BIGINT REFERENCES product_categories(id),
    description TEXT,
    image_url VARCHAR(500),
    purchase_unit VARCHAR(20), -- 箱
    stock_unit VARCHAR(20) NOT NULL, -- 瓶
    sale_unit VARCHAR(20) NOT NULL, -- 杯
    purchase_to_stock_ratio DECIMAL(10,4), -- 1箱=24瓶 -> 24
    stock_to_sale_ratio DECIMAL(10,4), -- 1瓶=5杯 -> 5
    cost_price DECIMAL(10,2), -- 成本价(库存单位)
    sale_price DECIMAL(10,2) NOT NULL, -- 售价(销售单位)
    shelf_life_days INT, -- 保质期天数
    storage_type VARCHAR(20) DEFAULT 'normal', -- normal/cold/frozen
    loss_rate DECIMAL(5,4) DEFAULT 0, -- 损耗率
    safety_stock DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active/inactive/deleted
    sort_order INT DEFAULT 0,
    attributes JSONB, -- 扩展属性
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_merchant ON products(merchant_id, status);

CREATE TABLE product_categories (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    name VARCHAR(100) NOT NULL,
    parent_id BIGINT REFERENCES product_categories(id),
    sort_order INT DEFAULT 0,
    icon_url VARCHAR(500)
);

-- BOM物料清单
CREATE TABLE product_bom (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id), -- 成品
    material_id BIGINT NOT NULL REFERENCES products(id), -- 原料
    quantity DECIMAL(10,4) NOT NULL, -- 用量(库存单位)
    unit VARCHAR(20) NOT NULL
);

-- ==================== 库存 ====================

CREATE TABLE inventory (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    current_stock DECIMAL(10,2) NOT NULL DEFAULT 0, -- 当前库存(库存单位)
    available_stock DECIMAL(10,2) NOT NULL DEFAULT 0, -- 可用库存
    locked_stock DECIMAL(10,2) NOT NULL DEFAULT 0, -- 锁定库存(已下单未核销)
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, product_id)
);

-- 库存批次(易腐品)
CREATE TABLE inventory_batches (
    id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL REFERENCES inventory(id),
    batch_no VARCHAR(50),
    quantity DECIMAL(10,2) NOT NULL,
    production_date DATE,
    expiry_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active/expired/consumed/damaged
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_batches_expiry ON inventory_batches(expiry_date, status);

-- 库存变动记录
CREATE TABLE inventory_logs (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    change_type VARCHAR(30) NOT NULL, -- purchase/sale/loss/adjustment/stocktake
    change_qty DECIMAL(10,2) NOT NULL,
    before_qty DECIMAL(10,2) NOT NULL,
    after_qty DECIMAL(10,2) NOT NULL,
    batch_id BIGINT REFERENCES inventory_batches(id),
    reference_id BIGINT, -- 关联订单/采购单ID
    remark TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- ==================== 订单 ====================

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) UNIQUE NOT NULL,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending_payment',
    -- pending_payment/paid/preparing/ready/completed/cancelled/refunded
    total_amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    pay_amount DECIMAL(10,2) NOT NULL,
    coupon_id BIGINT,
    pickup_time TIMESTAMPTZ,
    verify_code VARCHAR(10), -- 核销码
    verified_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancel_reason TEXT,
    wx_transaction_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_merchant_status ON orders(merchant_id, status);
CREATE INDEX idx_orders_user ON orders(user_id, created_at DESC);
CREATE INDEX idx_orders_verify ON orders(merchant_id, verify_code) WHERE status = 'paid';

CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    product_id BIGINT NOT NULL,
    product_name VARCHAR(200) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    spec VARCHAR(100)
);

-- ==================== 促销 ====================

CREATE TABLE promotions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    type VARCHAR(20) NOT NULL, -- clearance/flash_sale/new_arrival
    trigger_type VARCHAR(20), -- auto/manual
    trigger_reason VARCHAR(20), -- overstock/expiring
    original_price DECIMAL(10,2) NOT NULL,
    promo_price DECIMAL(10,2) NOT NULL,
    discount_rate DECIMAL(4,2), -- 0.7 = 7折
    predicted_sell_days INT,
    predicted_profit DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending/active/completed/cancelled
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    actual_sold_qty INT DEFAULT 0,
    actual_profit DECIMAL(10,2),
    alliance_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 优惠券
CREATE TABLE coupons (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL, -- fixed/percent
    value DECIMAL(10,2) NOT NULL, -- 金额或折扣
    min_amount DECIMAL(10,2) DEFAULT 0, -- 最低消费
    total_qty INT NOT NULL,
    used_qty INT DEFAULT 0,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE user_coupons (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    coupon_id BIGINT NOT NULL REFERENCES coupons(id),
    status VARCHAR(20) NOT NULL DEFAULT 'unused', -- unused/used/expired
    used_at TIMESTAMPTZ,
    order_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==================== 流量联盟 ====================

CREATE TABLE alliance_members (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    promo_credits INT DEFAULT 0, -- 推广券余额
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE alliance_exposures (
    id BIGSERIAL PRIMARY KEY,
    promotion_id BIGINT NOT NULL REFERENCES promotions(id),
    source_merchant_id BIGINT NOT NULL, -- 曝光在哪家店
    target_merchant_id BIGINT NOT NULL, -- 商品归属店
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    conversions INT DEFAULT 0,
    date DATE NOT NULL,
    UNIQUE(promotion_id, source_merchant_id, date)
);

-- ==================== AI预测 ====================

CREATE TABLE sales_predictions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    predict_date DATE NOT NULL,
    predicted_qty DECIMAL(10,2) NOT NULL,
    confidence_lower DECIMAL(10,2),
    confidence_upper DECIMAL(10,2),
    actual_qty DECIMAL(10,2), -- 事后回填
    model_version VARCHAR(20),
    features JSONB, -- 使用的特征
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, product_id, predict_date, model_version)
);

CREATE TABLE replenishment_suggestions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    current_stock DECIMAL(10,2) NOT NULL,
    predicted_daily_demand DECIMAL(10,2) NOT NULL,
    safety_stock DECIMAL(10,2) NOT NULL,
    suggested_qty DECIMAL(10,2) NOT NULL,
    expected_profit DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending/accepted/rejected
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==================== 计费 ====================

CREATE TABLE billing_records (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    plan VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending/paid/overdue
    paid_at TIMESTAMPTZ,
    wx_transaction_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==================== 通知与公告 ====================

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    target_type VARCHAR(20) NOT NULL, -- all/merchant/user
    target_id BIGINT, -- NULL=全体
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) NOT NULL, -- system/activity/alert
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==================== 行业模板 ====================

CREATE TABLE industry_templates (
    id BIGSERIAL PRIMARY KEY,
    industry VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    categories JSONB NOT NULL, -- 预设分类
    bom_templates JSONB, -- 预设BOM
    units_preset JSONB, -- 预设单位
    theme_config JSONB, -- 主题色/布局
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 4. API规范

### 4.1 设计原则

- RESTful风格，资源导向
- 版本化：`/api/v1/...`
- 统一响应格式：`{ "code": 0, "message": "success", "data": {} }`
- 分页：`?page=1&page_size=20`，响应包含 `total`
- 认证：JWT Bearer Token，Access Token (2h) + Refresh Token (7d)
- 错误码：业务错误 4xxxx，系统错误 5xxxx

### 4.2 认证方案

```
POST /api/v1/auth/login/phone     -- 手机号+验证码登录
POST /api/v1/auth/login/wechat    -- 微信授权登录
POST /api/v1/auth/refresh         -- 刷新Token
POST /api/v1/auth/logout          -- 登出
```

JWT Payload:
```json
{
  "uid": 12345,
  "role": "merchant",
  "mid": 100,         // merchant_id (商家角色)
  "exp": 1716200000
}
```

### 4.3 核心API列表

#### 平台管理端

```
-- 商家管理
GET    /api/v1/admin/merchants              -- 商家列表(筛选/分页)
GET    /api/v1/admin/merchants/:id          -- 商家详情
PUT    /api/v1/admin/merchants/:id/audit    -- 审核(通过/驳回)
PUT    /api/v1/admin/merchants/:id/freeze   -- 冻结/解冻
GET    /api/v1/admin/merchants/:id/stats    -- 商家经营数据

-- 模板管理
GET    /api/v1/admin/templates              -- 模板列表
POST   /api/v1/admin/templates              -- 创建模板
PUT    /api/v1/admin/templates/:id          -- 编辑模板

-- 小程序版本
POST   /api/v1/admin/miniapp/versions       -- 上传版本
POST   /api/v1/admin/miniapp/publish        -- 发布
POST   /api/v1/admin/miniapp/rollback       -- 回滚

-- 数据大盘
GET    /api/v1/admin/dashboard              -- 平台数据概览

-- 财务
GET    /api/v1/admin/billing                -- 账单列表
GET    /api/v1/admin/billing/plans          -- 计费方案
PUT    /api/v1/admin/billing/plans/:id      -- 修改方案

-- 运营活动
POST   /api/v1/admin/activities             -- 创建活动
GET    /api/v1/admin/activities             -- 活动列表
POST   /api/v1/admin/notifications          -- 推送通知
```

#### 商家端

```
-- 工作台
GET    /api/v1/merchant/dashboard           -- 今日概览
GET    /api/v1/merchant/alerts              -- 预警看板
GET    /api/v1/merchant/todos               -- 待办事项

-- 商品管理
GET    /api/v1/merchant/products            -- 商品列表
POST   /api/v1/merchant/products            -- 新增商品
PUT    /api/v1/merchant/products/:id        -- 编辑商品
DELETE /api/v1/merchant/products/:id        -- 删除商品
POST   /api/v1/merchant/products/voice      -- 语音录入
POST   /api/v1/merchant/products/import     -- 模板导入

-- BOM管理
GET    /api/v1/merchant/products/:id/bom    -- 获取BOM
PUT    /api/v1/merchant/products/:id/bom    -- 设置BOM

-- 库存
GET    /api/v1/merchant/inventory           -- 库存列表
POST   /api/v1/merchant/inventory/purchase  -- 入库
POST   /api/v1/merchant/inventory/stocktake -- 盘点
GET    /api/v1/merchant/inventory/batches   -- 批次列表

-- AI决策
GET    /api/v1/merchant/ai/replenishment    -- 补货建议
POST   /api/v1/merchant/ai/replenishment/confirm -- 确认补货
GET    /api/v1/merchant/ai/promotions       -- 促销建议
POST   /api/v1/merchant/ai/promotions/execute -- 执行促销
POST   /api/v1/merchant/ai/chat            -- AI顾问对话

-- 订单
GET    /api/v1/merchant/orders              -- 订单列表
GET    /api/v1/merchant/orders/:id          -- 订单详情
POST   /api/v1/merchant/orders/verify       -- 核销

-- 联盟
GET    /api/v1/merchant/alliance/status     -- 联盟状态
POST   /api/v1/merchant/alliance/join       -- 加入联盟
GET    /api/v1/merchant/alliance/stats      -- 联盟效果
POST   /api/v1/merchant/alliance/boost      -- 使用推广券

-- 店铺
GET    /api/v1/merchant/shop                -- 店铺信息
PUT    /api/v1/merchant/shop                -- 修改店铺
GET    /api/v1/merchant/subscription        -- 当前订阅
```

#### C端用户

```
-- 店铺浏览
GET    /api/v1/shop/:merchant_id            -- 店铺首页
GET    /api/v1/shop/:merchant_id/products   -- 商品列表
GET    /api/v1/shop/:merchant_id/products/:id -- 商品详情
GET    /api/v1/shop/:merchant_id/promotions -- 促销列表
GET    /api/v1/shop/:merchant_id/nearby     -- 附近推荐(联盟)

-- 购物车
GET    /api/v1/cart                         -- 获取购物车
POST   /api/v1/cart/items                   -- 添加商品
PUT    /api/v1/cart/items/:id               -- 修改数量
DELETE /api/v1/cart/items/:id               -- 删除商品

-- 订单
POST   /api/v1/orders                       -- 创建订单
GET    /api/v1/orders                       -- 订单列表
GET    /api/v1/orders/:id                   -- 订单详情
POST   /api/v1/orders/:id/pay              -- 支付
POST   /api/v1/orders/:id/cancel           -- 取消
POST   /api/v1/orders/:id/refund           -- 退款

-- 优惠券
GET    /api/v1/coupons/available            -- 可领优惠券
POST   /api/v1/coupons/:id/claim           -- 领取
GET    /api/v1/coupons/mine                 -- 我的优惠券
```

---

## 5. AI引擎设计

### 5.1 架构

```
┌───────────────┐       gRPC        ┌───────────────────┐
│  Go主服务     │  ◄──────────────►  │  AI引擎(Python)   │
│  (调度/缓存)  │                    │  FastAPI + gRPC   │
└───────────────┘                    └───────────────────┘
                                              │
                                     ┌────────┴────────┐
                                     │                 │
                               ┌─────▼─────┐   ┌──────▼──────┐
                               │ 预测模型   │   │  LLM服务    │
                               │ (Prophet/  │   │ (AI顾问)   │
                               │  XGBoost)  │   │            │
                               └───────────┘   └────────────┘
```

### 5.2 预测模型

| 模型 | 用途 | 输入特征 |
|:---|:---|:---|
| Prophet | 时间序列基准预测 | 历史日销量、星期效应、节假日 |
| XGBoost | 精细预测 | 天气、温度、促销、周边事件 |
| 规则引擎 | 冷启动/数据不足 | 行业基准、安全库存系数 |

### 5.3 数据流

1. **数据采集**：每日凌晨聚合前日销售数据
2. **特征工程**：拼接天气API、日历特征、促销状态
3. **模型训练**：每周离线重新训练(数据>30天的商户)
4. **实时推理**：商家打开工作台时触发补货/促销建议计算
5. **反馈闭环**：记录建议采纳率和实际效果，优化模型

### 5.4 AI顾问

- 基于大语言模型(Claude/GPT-4)
- 将商家数据摘要 + 问题 构成Prompt
- 流式返回、支持语音交互
- 上下文窗口保留近3轮对话

---

## 6. 性能方案

### 6.1 缓存策略

| 缓存层级 | 方案 | TTL |
|:---|:---|:---|
| CDN | 商品图片、静态资源 | 7天 |
| Redis L1 | 商品信息、店铺信息 | 5分钟 |
| Redis L2 | AI预测结果 | 1小时 |
| 本地缓存 | 配置信息、模板数据 | 10分钟 |

### 6.2 数据库优化

- 分区表：`inventory_logs` 按月分区
- 读写分离：主库写，从库读报表
- 连接池：PgBouncer (最大200连接)
- 索引：覆盖高频查询路径

### 6.3 高并发(秒杀)方案

```
用户请求 → Nginx限流(令牌桶)
         → Redis预扣库存(Lua原子操作)
         → 消息队列异步创建订单
         → 数据库最终一致
```

### 6.4 首屏优化

- 小程序分包加载
- 接口数据预加载(onLaunch)
- 骨架屏
- 图片懒加载 + WebP格式

---

## 7. 安全方案

### 7.1 认证授权

- JWT + RBAC权限模型
- 角色：super_admin / admin / merchant / staff / consumer
- 接口级权限校验中间件
- 敏感操作二次验证(短信)

### 7.2 数据安全

- 传输层：HTTPS/TLS 1.3
- 存储层：敏感字段AES-256加密(手机号、身份证)
- 商家数据逻辑隔离：所有查询强制带 merchant_id 条件
- 联盟数据聚合：≥5家数据方可展示
- SQL注入防护：参数化查询(GORM)
- XSS防护：输入过滤 + CSP头

### 7.3 审计日志

- 关键操作记录：操作人、时间、IP、操作内容
- 管理后台操作全量审计
- 日志保留180天

---

## 8. 项目目录结构

```
linximanager/
├── docs/                          # 文档
│   ├── architecture.md
│   └── api/
├── backend/                       # Go后端
│   ├── cmd/
│   │   └── server/main.go        # 入口
│   ├── internal/
│   │   ├── config/               # 配置
│   │   ├── middleware/           # 中间件(auth/cors/ratelimit/logger)
│   │   ├── handler/             # HTTP处理器
│   │   │   ├── admin/
│   │   │   ├── merchant/
│   │   │   └── consumer/
│   │   ├── service/             # 业务逻辑
│   │   ├── repository/          # 数据访问
│   │   ├── model/               # 数据模型
│   │   └── pkg/                 # 内部工具包
│   ├── migrations/              # 数据库迁移
│   ├── api/                     # API proto/OpenAPI定义
│   ├── go.mod
│   └── go.sum
├── ai-engine/                   # Python AI服务
│   ├── app/
│   │   ├── predictor/
│   │   ├── promoter/
│   │   └── advisor/
│   ├── models/                  # 训练好的模型
│   ├── requirements.txt
│   └── Dockerfile
├── frontend/                    # 前端
│   ├── admin/                   # 管理后台(Vue3 + Element Plus)
│   │   ├── src/
│   │   └── package.json
│   ├── merchant-app/            # 商家小程序(uni-app)
│   │   ├── src/
│   │   └── package.json
│   └── consumer-app/            # C端小程序(uni-app)
│       ├── src/
│       └── package.json
├── deploy/                      # 部署配置
│   ├── docker/
│   ├── k8s/
│   └── nginx/
├── .github/
│   └── workflows/               # CI/CD
└── docker-compose.yml
```

---

## 9. 部署架构概述

### 9.1 开发环境

```yaml
# docker-compose.yml 启动全部依赖
services:
  postgres, redis, minio, meilisearch, ai-engine, backend, admin-web
```

### 9.2 生产环境

```
                    ┌──── CDN (静态资源) ────┐
                    │                       │
用户 ─── DNS ─── [SLB/Nginx] ─── [K8s Ingress]
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
              [backend pods x3]  [ai-engine x2]   [admin-web]
                    │
         ┌─────────┼─────────┐
         │         │         │
    [PostgreSQL] [Redis]  [MinIO]
     (主从)      (哨兵)    (集群)
```

### 9.3 关键指标

- 后端实例：最少3副本，HPA自动扩缩(CPU>70%)
- 数据库：主从复制，每日自动备份
- Redis：哨兵模式，内存告警阈值80%
- 零停机部署：Rolling Update策略

---

## 10. 开发规范

### 10.1 Git规范

- 分支：main / develop / feature/* / hotfix/*
- Commit：Conventional Commits (feat/fix/docs/chore)
- PR：必须Code Review + CI通过

### 10.2 API版本策略

- URL版本化：/api/v1/...
- 大版本变更时v2与v1并行运行，给出迁移窗口

### 10.3 错误处理

```go
// 统一错误响应
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 业务错误码
// 40001 - 参数校验失败
// 40100 - 未认证
// 40300 - 无权限
// 40400 - 资源不存在
// 42900 - 请求过于频繁
// 50000 - 内部错误
```
