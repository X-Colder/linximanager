---
name: 灵犀掌柜前端项目结构
description: 三端前端项目的技术栈、目录结构和核心组件
type: project
---

三端前端代码已在 `/Users/yaojun72/Documents/workspace/llm/linximanager/frontend/` 下完整创建。

**Why:** 2026-05-20 team-lead 分配前端开发任务，基于架构和UI设计文档实现。

**How to apply:** 后续开发时直接定位到对应子目录，遵循已有代码结构扩展。

## admin（管理后台）
- 技术栈：Vue3 + TypeScript + Element Plus + Vite + Pinia + ECharts
- Axios 封装：`src/api/request.ts`，含 JWT 拦截器 + Refresh Token 自动刷新（并发锁）
- auth store：`src/stores/auth.ts`，token 持久化到 localStorage
- 路由守卫：`src/router/index.ts`，未登录跳转 /login
- 已实现页面：登录、数据大盘（ECharts）、商家列表（表格+筛选+分页）、入驻审核（左右分栏）、模板管理

## merchant-app（商家小程序）
- 技术栈：Vue3 + TypeScript + uni-app + Pinia
- API 封装：`src/api/request.ts`，基于 `uni.request`，含 loading/error/401跳转
- 5个 tabBar 页面：首页工作台/商品库存/AI决策/订单管理/我的
- 核心组件：`AlertCard.vue`（红黄蓝三色预警）、`InventoryShelf.vue`（货架+列表双模式）、`PromoCompare.vue`（方案对比）

## consumer-app（C端小程序）
- 技术栈：Vue3 + TypeScript + uni-app
- 3个 tabBar 页面：店铺首页/购物车/我的订单
- 核心组件：`ProductCard.vue`（横向布局+促销标签）、`CouponCard.vue`（锯齿优惠券）
- 购物车本地持久化 `uni.setStorageSync('cart', ...)`

## API 对接规范
- 统一响应格式：`{ code: 0, message: "success", data: T }`
- 错误码：0=成功，401=未认证，403=无权限
- 分页：`?page=1&page_size=20`，返回 `{ list, total, page, page_size }`
