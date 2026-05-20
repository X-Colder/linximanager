-- 灵犀掌柜 初始 Schema
-- 适用数据库：PostgreSQL 16

BEGIN;

-- ==================== 扩展 ====================
CREATE EXTENSION IF NOT EXISTS postgis;

-- ==================== 用户与认证 ====================

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'consumer',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    openid VARCHAR(100) UNIQUE,
    unionid VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- ==================== 商家 ====================

CREATE TABLE IF NOT EXISTS merchants (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    shop_name VARCHAR(200) NOT NULL,
    industry VARCHAR(50) NOT NULL,
    logo_url VARCHAR(500),
    address TEXT,
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    contact_phone VARCHAR(20),
    business_hours VARCHAR(100),
    announcement TEXT,
    license_url VARCHAR(500),
    id_card_url VARCHAR(500),
    shop_photos JSONB,
    audit_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    audit_remark TEXT,
    version_plan VARCHAR(20) NOT NULL DEFAULT 'basic',
    plan_expire_at TIMESTAMPTZ,
    miniapp_appid VARCHAR(100),
    miniapp_status VARCHAR(20) DEFAULT 'pending',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchants_industry ON merchants(industry);
CREATE INDEX IF NOT EXISTS idx_merchants_status ON merchants(status);
CREATE INDEX IF NOT EXISTS idx_merchants_audit_status ON merchants(audit_status);
CREATE INDEX IF NOT EXISTS idx_merchants_user_id ON merchants(user_id);

-- ==================== 商品分类 ====================

CREATE TABLE IF NOT EXISTS product_categories (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    name VARCHAR(100) NOT NULL,
    parent_id BIGINT REFERENCES product_categories(id),
    sort_order INT DEFAULT 0,
    icon_url VARCHAR(500)
);

CREATE INDEX IF NOT EXISTS idx_product_categories_merchant ON product_categories(merchant_id);

-- ==================== 商品 ====================

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    name VARCHAR(200) NOT NULL,
    category_id BIGINT REFERENCES product_categories(id),
    description TEXT,
    image_url VARCHAR(500),
    purchase_unit VARCHAR(20),
    stock_unit VARCHAR(20) NOT NULL,
    sale_unit VARCHAR(20) NOT NULL,
    purchase_to_stock_ratio DECIMAL(10,4),
    stock_to_sale_ratio DECIMAL(10,4),
    cost_price DECIMAL(10,2),
    sale_price DECIMAL(10,2) NOT NULL,
    shelf_life_days INT,
    storage_type VARCHAR(20) DEFAULT 'normal',
    loss_rate DECIMAL(5,4) DEFAULT 0,
    safety_stock DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    sort_order INT DEFAULT 0,
    attributes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_merchant ON products(merchant_id, status);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);

-- BOM 物料清单
CREATE TABLE IF NOT EXISTS product_bom (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    material_id BIGINT NOT NULL REFERENCES products(id),
    quantity DECIMAL(10,4) NOT NULL,
    unit VARCHAR(20) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_product_bom_product ON product_bom(product_id);

-- ==================== 库存 ====================

CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    current_stock DECIMAL(10,2) NOT NULL DEFAULT 0,
    available_stock DECIMAL(10,2) NOT NULL DEFAULT 0,
    locked_stock DECIMAL(10,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_merchant ON inventory(merchant_id);

-- 库存批次
CREATE TABLE IF NOT EXISTS inventory_batches (
    id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL REFERENCES inventory(id),
    batch_no VARCHAR(50),
    quantity DECIMAL(10,2) NOT NULL,
    production_date DATE,
    expiry_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_batches_expiry ON inventory_batches(expiry_date, status);

-- 库存变动记录（按月分区）
CREATE TABLE IF NOT EXISTS inventory_logs (
    id BIGSERIAL,
    merchant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    change_type VARCHAR(30) NOT NULL,
    change_qty DECIMAL(10,2) NOT NULL,
    before_qty DECIMAL(10,2) NOT NULL,
    after_qty DECIMAL(10,2) NOT NULL,
    batch_id BIGINT REFERENCES inventory_batches(id),
    reference_id BIGINT,
    remark TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- 创建分区（2026年）
CREATE TABLE IF NOT EXISTS inventory_logs_2026_01 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_02 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_03 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_04 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_05 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_06 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_07 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_08 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_09 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_10 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_11 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE IF NOT EXISTS inventory_logs_2026_12 PARTITION OF inventory_logs
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

CREATE INDEX IF NOT EXISTS idx_inventory_logs_merchant ON inventory_logs(merchant_id, created_at);

-- ==================== 订单 ====================

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) UNIQUE NOT NULL,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending_payment',
    total_amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    pay_amount DECIMAL(10,2) NOT NULL,
    coupon_id BIGINT,
    pickup_time TIMESTAMPTZ,
    verify_code VARCHAR(10),
    verified_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancel_reason TEXT,
    wx_transaction_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_merchant_status ON orders(merchant_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_verify ON orders(merchant_id, verify_code) WHERE status = 'paid';

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL,
    product_name VARCHAR(200) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    spec VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);

-- ==================== 促销 ====================

CREATE TABLE IF NOT EXISTS promotions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    type VARCHAR(20) NOT NULL,
    trigger_type VARCHAR(20),
    trigger_reason VARCHAR(20),
    original_price DECIMAL(10,2) NOT NULL,
    promo_price DECIMAL(10,2) NOT NULL,
    discount_rate DECIMAL(4,2),
    predicted_sell_days INT,
    predicted_profit DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    actual_sold_qty INT DEFAULT 0,
    actual_profit DECIMAL(10,2),
    alliance_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_promotions_merchant ON promotions(merchant_id, status);

-- 优惠券
CREATE TABLE IF NOT EXISTS coupons (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    min_amount DECIMAL(10,2) DEFAULT 0,
    total_qty INT NOT NULL,
    used_qty INT DEFAULT 0,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS user_coupons (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    coupon_id BIGINT NOT NULL REFERENCES coupons(id),
    status VARCHAR(20) NOT NULL DEFAULT 'unused',
    used_at TIMESTAMPTZ,
    order_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_coupons_user ON user_coupons(user_id);
CREATE INDEX IF NOT EXISTS idx_user_coupons_coupon ON user_coupons(coupon_id);

-- ==================== 流量联盟 ====================

CREATE TABLE IF NOT EXISTS alliance_members (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    promo_credits INT DEFAULT 0,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id)
);

CREATE TABLE IF NOT EXISTS alliance_exposures (
    id BIGSERIAL PRIMARY KEY,
    promotion_id BIGINT NOT NULL REFERENCES promotions(id),
    source_merchant_id BIGINT NOT NULL,
    target_merchant_id BIGINT NOT NULL,
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    conversions INT DEFAULT 0,
    date DATE NOT NULL,
    UNIQUE(promotion_id, source_merchant_id, date)
);

-- ==================== AI预测 ====================

CREATE TABLE IF NOT EXISTS sales_predictions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    predict_date DATE NOT NULL,
    predicted_qty DECIMAL(10,2) NOT NULL,
    confidence_lower DECIMAL(10,2),
    confidence_upper DECIMAL(10,2),
    actual_qty DECIMAL(10,2),
    model_version VARCHAR(20),
    features JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, product_id, predict_date, model_version)
);

CREATE TABLE IF NOT EXISTS replenishment_suggestions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    current_stock DECIMAL(10,2) NOT NULL,
    predicted_daily_demand DECIMAL(10,2) NOT NULL,
    safety_stock DECIMAL(10,2) NOT NULL,
    suggested_qty DECIMAL(10,2) NOT NULL,
    expected_profit DECIMAL(10,2),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_replenishment_merchant ON replenishment_suggestions(merchant_id, status);

-- ==================== 计费 ====================

CREATE TABLE IF NOT EXISTS billing_records (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    plan VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMPTZ,
    wx_transaction_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_billing_merchant ON billing_records(merchant_id);

-- ==================== 通知 ====================

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    target_type VARCHAR(20) NOT NULL,
    target_id BIGINT,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_target ON notifications(target_type, target_id);

-- ==================== 行业模板 ====================

CREATE TABLE IF NOT EXISTS industry_templates (
    id BIGSERIAL PRIMARY KEY,
    industry VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    categories JSONB NOT NULL,
    bom_templates JSONB,
    units_preset JSONB,
    theme_config JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_industry_templates_industry ON industry_templates(industry, status);

COMMIT;
