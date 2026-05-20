<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="navbar">
      <text class="navbar__title">灵犀掌柜</text>
      <text class="navbar__date">{{ today }}</text>
    </view>

    <scroll-view scroll-y class="scroll-content" @refresherrefresh="onRefresh" :refresher-triggered="refreshing" refresher-enabled>

      <!-- 今日概览 -->
      <view class="section-card">
        <view class="section-title">今日概览</view>
        <view class="overview-grid" v-if="store.dashboard">
          <view class="overview-item">
            <text class="overview-value">¥{{ store.dashboard.today_sales.toLocaleString() }}</text>
            <text class="overview-label">销售额</text>
            <text class="overview-growth" :class="store.dashboard.sales_growth >= 0 ? 'growth--up' : 'growth--down'">
              {{ store.dashboard.sales_growth >= 0 ? '↑' : '↓' }}{{ Math.abs(store.dashboard.sales_growth) }}%
            </text>
          </view>
          <view class="overview-item">
            <text class="overview-value">{{ store.dashboard.today_orders }}</text>
            <text class="overview-label">订单数</text>
            <text class="overview-growth" :class="store.dashboard.orders_growth >= 0 ? 'growth--up' : 'growth--down'">
              {{ store.dashboard.orders_growth >= 0 ? '↑' : '↓' }}{{ Math.abs(store.dashboard.orders_growth) }}%
            </text>
          </view>
          <view class="overview-item">
            <text class="overview-value">¥{{ store.dashboard.today_profit }}</text>
            <text class="overview-label">毛利</text>
            <text class="overview-growth" :class="store.dashboard.profit_growth >= 0 ? 'growth--up' : 'growth--down'">
              {{ store.dashboard.profit_growth >= 0 ? '↑' : '↓' }}{{ Math.abs(store.dashboard.profit_growth) }}%
            </text>
          </view>
        </view>
        <!-- 骨架屏 -->
        <view v-else class="overview-skeleton">
          <view v-for="i in 3" :key="i" class="skeleton-item" />
        </view>
      </view>

      <!-- 预警看板 -->
      <view class="section-card">
        <view class="section-header">
          <text class="section-title">预警看板</text>
          <text class="section-more" @tap="goToProducts">查看全部 ›</text>
        </view>

        <view v-if="alertGroups.red.length > 0">
          <view class="alert-group-title">🔴 缺货预警 ({{ alertGroups.red.length }})</view>
          <AlertCard
            v-for="alert in alertGroups.red"
            :key="alert.id"
            level="red"
            :title="alert.product_name"
            :subtitle="alert.subtitle"
            :description="alert.description"
            :action-text="alert.action_text"
            @action="handleAlertAction(alert)"
          />
        </view>

        <view v-if="alertGroups.yellow.length > 0">
          <view class="alert-group-title">🟡 积压预警 ({{ alertGroups.yellow.length }})</view>
          <AlertCard
            v-for="alert in alertGroups.yellow"
            :key="alert.id"
            level="yellow"
            :title="alert.product_name"
            :subtitle="alert.subtitle"
            :description="alert.description"
            :action-text="alert.action_text"
            @action="handleAlertAction(alert)"
          />
        </view>

        <view v-if="alertGroups.blue.length > 0">
          <view class="alert-group-title">🔵 临期提醒 ({{ alertGroups.blue.length }})</view>
          <AlertCard
            v-for="alert in alertGroups.blue"
            :key="alert.id"
            level="blue"
            :title="alert.product_name"
            :subtitle="alert.subtitle"
            :description="alert.description"
            :action-text="alert.action_text"
            @action="handleAlertAction(alert)"
          />
        </view>

        <view v-if="!hasAlerts && !store.loading" class="empty-alerts">
          <text class="empty-icon">✅</text>
          <text class="empty-text">一切正常，继续保持！</text>
        </view>
      </view>

      <!-- 待办事项 -->
      <view class="section-card">
        <view class="section-header">
          <text class="section-title">待办事项</text>
          <text class="todo-count">({{ store.todos.length }})</text>
        </view>
        <view class="todo-list">
          <view
            v-for="todo in store.todos"
            :key="todo.id"
            class="todo-item"
            :class="{ 'todo-item--done': todo.done }"
          >
            <view class="todo-dot" :class="todo.done ? 'todo-dot--done' : ''" />
            <text class="todo-text">{{ todo.content }}</text>
          </view>
          <view v-if="!store.todos.length && !store.loading" class="empty-state">
            <text class="empty-text">暂无待办事项</text>
          </view>
        </view>
      </view>

    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import AlertCard from '@/components/AlertCard.vue'
import { useMerchantStore, type Alert } from '@/stores/merchant'

const store = useMerchantStore()
const refreshing = ref(false)

const today = computed(() => {
  const d = new Date()
  const days = ['日', '一', '二', '三', '四', '五', '六']
  return `${d.getMonth() + 1}月${d.getDate()}日 周${days[d.getDay()]}`
})

const alertGroups = computed(() => ({
  red:    store.alerts.filter(a => a.level === 'red'),
  yellow: store.alerts.filter(a => a.level === 'yellow'),
  blue:   store.alerts.filter(a => a.level === 'blue')
}))

const hasAlerts = computed(() =>
  store.alerts.length > 0
)

async function loadData() {
  await Promise.all([
    store.fetchDashboard(),
    store.fetchAlerts(),
    store.fetchTodos()
  ])
}

async function onRefresh() {
  refreshing.value = true
  await loadData()
  refreshing.value = false
}

function handleAlertAction(alert: Alert) {
  if (alert.action_type === 'replenish') {
    uni.navigateTo({ url: '/pages/decision/index?tab=replenishment' })
  } else if (alert.action_type === 'promo' || alert.action_type === 'clearance') {
    uni.navigateTo({ url: `/pages/decision/index?tab=promo&productId=${alert.product_id}` })
  }
}

function goToProducts() {
  uni.switchTab({ url: '/pages/products/index' })
}

onMounted(() => loadData())
onShow(() => {
  // 每次进入页面刷新数据
  loadData()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f9fafb;
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 32rpx;
  background: #4a6cf7;
}

.navbar__title {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
}

.navbar__date {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.scroll-content {
  height: calc(100vh - 100rpx);
  padding: 24rpx 24rpx 40rpx;
}

.section-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 20rpx;
}

.section-header .section-title {
  margin-bottom: 0;
}

.section-more {
  font-size: 24rpx;
  color: #4a6cf7;
}

.todo-count { font-size: 24rpx; color: #9ca3af; }

/* 今日概览 */
.overview-grid {
  display: flex;
  gap: 0;
}

.overview-item {
  flex: 1;
  text-align: center;
  border-right: 1rpx solid #f3f4f6;
  padding: 0 16rpx;
}

.overview-item:last-child { border-right: none; }

.overview-value {
  display: block;
  font-size: 40rpx;
  font-weight: 700;
  color: #1f2937;
}

.overview-label {
  display: block;
  font-size: 24rpx;
  color: #9ca3af;
  margin-top: 4rpx;
}

.overview-growth {
  display: block;
  font-size: 22rpx;
  margin-top: 4rpx;
}

.growth--up   { color: #10b981; }
.growth--down { color: #ef4444; }

.overview-skeleton {
  display: flex;
  gap: 16rpx;
}

.skeleton-item {
  flex: 1;
  height: 80rpx;
  background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
  border-radius: 8rpx;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

/* 预警 */
.alert-group-title {
  font-size: 26rpx;
  font-weight: 600;
  color: #4b5563;
  margin-bottom: 12rpx;
  margin-top: 4rpx;
}

.empty-alerts {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40rpx 0;
  gap: 12rpx;
}

.empty-icon { font-size: 64rpx; }
.empty-text { font-size: 26rpx; color: #9ca3af; }

/* 待办 */
.todo-list { display: flex; flex-direction: column; gap: 16rpx; }

.todo-item {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.todo-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  border: 2rpx solid #9ca3af;
  flex-shrink: 0;
}

.todo-dot--done {
  background: #10b981;
  border-color: #10b981;
}

.todo-text {
  font-size: 28rpx;
  color: #1f2937;
}

.todo-item--done .todo-text {
  color: #9ca3af;
  text-decoration: line-through;
}

.empty-state { text-align: center; padding: 20rpx 0; }
</style>
