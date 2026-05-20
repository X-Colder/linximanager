# 灵犀掌柜前端测试计划

**文档日期**: 2026-05-20
**版本**: v1.0
**负责人**: qas（自动化测试 Agent）

---

## 一、测试范围概述

| 应用 | 技术栈 | 测试框架推荐 |
|------|--------|------------|
| admin（管理后台）| Vue3 + Element Plus + Vite | Vitest + Vue Test Utils + Playwright |
| merchant-app（商家小程序）| uni-app + Vue3 | Vitest + Vue Test Utils |
| consumer-app（C端小程序）| uni-app + Vue3 | Vitest + Vue Test Utils |

---

## 二、E2E 测试用例设计

### 2.1 管理后台登录流程（Playwright）

**测试文件**: `tests/e2e/admin/auth.spec.ts`

#### TC-E2E-001: 正常登录
- **Given**: 用户访问 `/login`，输入合法 admin 账号和密码
- **When**: 点击"登录"按钮
- **Then**: 跳转到 `/dashboard`，顶部显示用户名

#### TC-E2E-002: 错误密码登录
- **Given**: 用户输入正确手机号、错误密码
- **When**: 点击"登录"
- **Then**: 页面显示错误提示，不跳转

#### TC-E2E-003: 空表单提交
- **Given**: 手机号和密码均为空
- **When**: 点击"登录"
- **Then**: 表单校验提示出现，接口不被调用

#### TC-E2E-004: 登录后刷新保持登录态
- **Given**: 已成功登录
- **When**: 刷新页面
- **Then**: 保持登录态，不重定向到登录页

#### TC-E2E-005: 无权限页面跳转
- **Given**: 未登录用户直接访问 `/merchant-list`
- **When**: 路由守卫检测
- **Then**: 被重定向到 `/login`

```typescript
// 示例实现思路（Playwright）
import { test, expect } from '@playwright/test';

test('管理员正常登录', async ({ page }) => {
  await page.goto('/login');
  await page.fill('[data-testid="phone-input"]', '13800138000');
  await page.fill('[data-testid="code-input"]', '123456');
  await page.click('[data-testid="login-btn"]');
  await expect(page).toHaveURL('/dashboard');
});
```

---

### 2.2 商家端核心操作流程（Playwright/uni-app）

**测试文件**: `tests/e2e/merchant/inventory.spec.ts`

#### TC-E2E-010: 商品入库操作
- **Given**: 商家已登录，进入库存管理页
- **When**: 选择商品、输入入库数量 100，提交
- **Then**: 库存列表中该商品数量增加 100，入库记录可查

#### TC-E2E-011: 库存盘点操作
- **Given**: 商品当前库存为 50
- **When**: 输入盘点数量 45，提交
- **Then**: 库存变为 45，变动记录显示 -5（类型: stocktake）

#### TC-E2E-012: 订单核销流程
- **Given**: 存在一笔 status=paid 的订单，有核销码
- **When**: 商家输入 6 位核销码点击核销
- **Then**: 订单状态变为 completed，库存对应减少

#### TC-E2E-013: 预警看板展示
- **Given**: 某商品库存低于安全库存
- **When**: 商家进入工作台
- **Then**: 预警看板显示该商品缺货预警

#### TC-E2E-014: AI 补货建议
- **Given**: 商家进入 AI 决策页面
- **When**: 系统返回补货建议列表
- **Then**: 显示商品名称、建议补货量、预期收益

---

### 2.3 C 端下单流程（Playwright）

**测试文件**: `tests/e2e/consumer/order.spec.ts`

#### TC-E2E-020: 完整下单流程
- **Given**: 用户已登录，进入店铺商品页
- **When**: 添加商品到购物车 → 结算 → 确认订单 → 模拟支付
- **Then**: 订单创建成功，状态为 paid，购物车清空

#### TC-E2E-021: 购物车数量修改
- **Given**: 购物车中有商品 A 数量为 1
- **When**: 点击 + 号增加到 3
- **Then**: 购物车数量变为 3，合计金额同步更新

#### TC-E2E-022: 库存不足提示
- **Given**: 商品 A 库存为 2
- **When**: 用户尝试加入购物车 5 个
- **Then**: 提示"库存不足"，购物车数量上限为 2

#### TC-E2E-023: 取消待支付订单
- **Given**: 存在一笔 pending_payment 订单
- **When**: 用户点击取消
- **Then**: 订单状态变为 cancelled，锁定库存释放

#### TC-E2E-024: 优惠券使用
- **Given**: 用户持有满 50 减 10 的优惠券，订单总额 60
- **When**: 选择使用该优惠券下单
- **Then**: 实付金额为 50，优惠金额显示 -10

---

## 三、组件单元测试清单

### 3.1 AlertCard.vue（预警卡片）

**测试文件**: `tests/unit/merchant/AlertCard.test.ts`

| 测试用例 | 场景 | 期望行为 |
|---------|------|---------|
| TC-C-001 | 传入 type=danger, message="库存不足" | 渲染红色背景，显示图标和消息文本 |
| TC-C-002 | 传入 type=warning | 渲染橙色背景 |
| TC-C-003 | 传入 type=info | 渲染蓝色背景 |
| TC-C-004 | 点击"查看详情"按钮 | 触发 @click-detail 事件 |
| TC-C-005 | message 为空字符串 | 不显示消息区域或显示默认占位符 |
| TC-C-006 | 传入 count 数字 | 显示预警数量徽章 |

```typescript
// 示例（Vitest + Vue Test Utils）
import { mount } from '@vue/test-utils';
import AlertCard from '@/components/AlertCard.vue';

describe('AlertCard', () => {
  it('渲染 danger 类型', () => {
    const wrapper = mount(AlertCard, {
      props: { type: 'danger', message: '库存不足' }
    });
    expect(wrapper.classes()).toContain('alert-danger');
    expect(wrapper.text()).toContain('库存不足');
  });
});
```

---

### 3.2 InventoryShelf.vue（库存货架）

**测试文件**: `tests/unit/merchant/InventoryShelf.test.ts`

| 测试用例 | 场景 | 期望行为 |
|---------|------|---------|
| TC-C-010 | 传入正常库存列表 | 正确渲染每行商品名称和库存量 |
| TC-C-011 | 库存量低于安全库存 | 该行显示红色警示 |
| TC-C-012 | 传入空列表 | 渲染空状态提示（"暂无库存数据"） |
| TC-C-013 | 点击"入库"按钮 | 触发 @purchase 事件并携带 productId |
| TC-C-014 | 点击"盘点"按钮 | 触发 @stocktake 事件 |
| TC-C-015 | 数据加载中状态 | 显示骨架屏或 Loading |

---

### 3.3 PromoCompare.vue（促销对比）

**测试文件**: `tests/unit/merchant/PromoCompare.test.ts`

| 测试用例 | 场景 | 期望行为 |
|---------|------|---------|
| TC-C-020 | 传入原价和促销价 | 显示折扣率计算正确（如7折） |
| TC-C-021 | 预测销售天数 > 保质期 | 显示红色警告 |
| TC-C-022 | 点击"执行促销"按钮 | 触发 @execute 事件 |
| TC-C-023 | 点击"忽略"按钮 | 触发 @dismiss 事件 |
| TC-C-024 | 折扣率 = 1.0（不打折）| 不显示折扣徽章 |
| TC-C-025 | 预期收益为负 | 显示亏损警告样式 |

---

### 3.4 ProductCard.vue（C端商品卡片）

**测试文件**: `tests/unit/consumer/ProductCard.test.ts`

| 测试用例 | 场景 | 期望行为 |
|---------|------|---------|
| TC-C-030 | 传入商品基本信息 | 显示商品名称、价格、图片 |
| TC-C-031 | 商品有促销价 | 原价显示删除线，促销价高亮 |
| TC-C-032 | 库存为 0 | 按钮显示"售罄"且不可点击 |
| TC-C-033 | 点击"加入购物车" | 触发 @add-to-cart 事件 |
| TC-C-034 | 图片加载失败 | 显示默认占位图 |
| TC-C-035 | 商品名称超长（> 20字） | 文本截断并显示省略号 |
| TC-C-036 | 传入 promotion 标签 | 显示促销角标 |

---

## 四、性能测试基准

### 4.1 前端性能指标

| 指标 | 目标值 | 测量方法 |
|------|-------|---------|
| 首屏加载时间（FCP）| ≤ 2s（4G网络）| Lighthouse / Playwright performance API |
| 最大内容绘制（LCP）| ≤ 2.5s | Lighthouse |
| 交互延迟（TTI）| ≤ 3s | Lighthouse |
| 包体积（JS bundle）| ≤ 500KB（gzip前）| Vite build 分析 |
| 小程序首屏渲染 | ≤ 1.5s | 微信开发者工具 Audits |

### 4.2 API 响应性能指标

| 接口类型 | P50 目标 | P99 目标 |
|---------|---------|---------|
| 登录接口 | ≤ 200ms | ≤ 500ms |
| 商品列表 | ≤ 150ms | ≤ 400ms |
| 库存查询 | ≤ 100ms | ≤ 300ms |
| 创建订单 | ≤ 300ms | ≤ 800ms |
| AI 建议（补货/促销）| ≤ 1s | ≤ 3s |
| AI 顾问对话（首字节）| ≤ 500ms | ≤ 1s |

### 4.3 性能测试工具推荐

```bash
# 前端性能
npx lighthouse http://localhost:5173 --output html

# API 压测（k6）
k6 run --vus 50 --duration 30s tests/perf/login.js

# 前端包体积分析
npx vite-bundle-visualizer
```

### 4.4 k6 压测脚本示例

```javascript
// tests/perf/login.js
import http from 'k6/http';
import { check } from 'k6';

export default function() {
  const res = http.post('http://localhost:8080/api/v1/auth/login/phone', JSON.stringify({
    phone: '13812345678',
    code: '123456'
  }), { headers: { 'Content-Type': 'application/json' } });

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
}
```

---

## 五、测试环境配置

### 5.1 Vitest 配置（管理后台）

```typescript
// frontend/admin/vite.config.ts（测试部分）
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      reporter: ['text', 'lcov'],
      exclude: ['node_modules/', 'src/router/'],
    },
  },
});
```

### 5.2 依赖安装

```bash
# 管理后台单元测试
cd frontend/admin
npm install -D vitest @vue/test-utils @vitejs/plugin-vue jsdom

# E2E 测试
npm install -D @playwright/test
npx playwright install
```

---

## 六、测试优先级矩阵

| 模块 | 优先级 | 测试类型 | 原因 |
|------|-------|---------|------|
| 用户登录/鉴权 | P0 | E2E + 单元 | 所有功能入口 |
| 库存入库/扣减 | P0 | E2E + 单元 | 核心业务，涉及金钱 |
| 订单创建/核销 | P0 | E2E | 最高频用户流 |
| AlertCard 预警 | P1 | 单元 | 商家决策依赖 |
| InventoryShelf | P1 | 单元 | 高频使用组件 |
| ProductCard | P1 | 单元 | C端核心展示 |
| PromoCompare | P2 | 单元 | AI功能辅助 |
| 优惠券流程 | P2 | E2E | 非核心路径 |
