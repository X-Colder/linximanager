<template>
  <view class="page">
    <!-- Tab 切换 -->
    <view class="tab-bar">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @tap="currentTab = tab.value"
      >
        {{ tab.label }}
      </view>
    </view>

    <scroll-view scroll-y class="order-list">
      <view v-for="order in filteredOrders" :key="order.id" class="order-card">
        <view class="order-header">
          <text class="order-shop">{{ order.shop_name }}</text>
          <text class="order-status" :class="`status--${order.status}`">
            {{ STATUS_LABEL[order.status] }}
          </text>
        </view>

        <view class="order-items">
          <view v-for="item in order.items" :key="item.product_id" class="order-item">
            <image class="item-img" :src="item.image_url || ''" mode="aspectFill" />
            <view class="item-detail">
              <text class="item-name">{{ item.product_name }}</text>
              <text class="item-qty">x{{ item.quantity }}</text>
            </view>
            <text class="item-price">¥{{ item.total_price }}</text>
          </view>
        </view>

        <view class="order-footer">
          <text class="order-total">实付 <text class="total-amount">¥{{ order.pay_amount }}</text></text>
        </view>

        <!-- 核销码（已付款/待自提） -->
        <view v-if="order.status === 'paid' || order.status === 'ready'" class="verify-code-section">
          <text class="verify-label">出示核销码</text>
          <view class="verify-code">{{ order.verify_code }}</view>
        </view>

        <!-- 操作按钮 -->
        <view class="order-actions">
          <button v-if="order.status === 'pending_payment'" class="action-btn action-btn--primary" @tap="payOrder(order)">
            去付款
          </button>
          <button v-if="order.status === 'pending_payment'" class="action-btn" @tap="cancelOrder(order)">
            取消订单
          </button>
          <button v-if="order.status === 'completed'" class="action-btn" @tap="reorder(order)">
            再来一单
          </button>
          <button v-if="order.status === 'completed'" class="action-btn action-btn--refund" @tap="refundOrder(order)">
            申请退款
          </button>
        </view>
      </view>

      <view v-if="!filteredOrders.length && !loading" class="empty-state">
        <text class="empty-icon">📋</text>
        <text class="empty-text">暂无订单</text>
        <button class="go-shop-btn" @tap="goShopping">去逛逛</button>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import http from '@/api/request'

type OrderStatus = 'pending_payment' | 'paid' | 'preparing' | 'ready' | 'completed' | 'cancelled' | 'refunded'

interface OrderItem {
  product_id: number
  product_name: string
  image_url?: string
  quantity: number
  total_price: number
}

interface Order {
  id: number
  shop_name: string
  status: OrderStatus
  pay_amount: number
  verify_code?: string
  items: OrderItem[]
  created_at: string
}

const STATUS_LABEL: Record<OrderStatus, string> = {
  pending_payment: '待付款',
  paid: '待备货',
  preparing: '备货中',
  ready: '待自提',
  completed: '已完成',
  cancelled: '已取消',
  refunded: '已退款'
}

const tabs = [
  { label: '全部', value: '' },
  { label: '待付款', value: 'pending_payment' },
  { label: '待自提', value: 'ready' },
  { label: '已完成', value: 'completed' }
]

const currentTab = ref('')
const orders = ref<Order[]>([])
const loading = ref(false)

const filteredOrders = computed(() => {
  if (!currentTab.value) return orders.value
  return orders.value.filter(o => o.status === currentTab.value)
})

async function payOrder(order: Order) {
  try {
    await http.post(`/orders/${order.id}/pay`, {}, true)
    order.status = 'paid'
    uni.showToast({ title: '支付成功', icon: 'success' })
  } catch {}
}

async function cancelOrder(order: Order) {
  await uni.showModal({ title: '提示', content: '确定要取消订单吗？' })
  try {
    await http.post(`/orders/${order.id}/cancel`, {}, true)
    order.status = 'cancelled'
    uni.showToast({ title: '已取消', icon: 'success' })
  } catch {}
}

async function refundOrder(order: Order) {
  await uni.showModal({ title: '申请退款', content: '确定要申请退款吗？' })
  try {
    await http.post(`/orders/${order.id}/refund`, {}, true)
    order.status = 'refunded'
    uni.showToast({ title: '退款申请已提交', icon: 'success' })
  } catch {}
}

function reorder(order: Order) {
  uni.switchTab({ url: '/pages/index/index' })
}

function goShopping() {
  uni.switchTab({ url: '/pages/index/index' })
}

onMounted(async () => {
  loading.value = true
  try {
    orders.value = await http.get<Order[]>('/orders')
  } catch {
    orders.value = [
      {
        id: 1,
        shop_name: '老王烧烤',
        status: 'paid',
        pay_amount: 74,
        verify_code: 'LX8821',
        items: [
          { product_id: 1, product_name: '烤鸡翅(5串)', quantity: 2, total_price: 30 },
          { product_id: 3, product_name: '青岛啤酒', quantity: 4, total_price: 24 }
        ],
        created_at: new Date().toISOString()
      }
    ]
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; }

.tab-bar {
  display: flex;
  background: #fff;
  border-bottom: 1rpx solid #f3f4f6;
  overflow-x: auto;
}

.tab-item {
  flex: 1;
  min-width: 120rpx;
  text-align: center;
  padding: 24rpx 0;
  font-size: 26rpx;
  color: #9ca3af;
  border-bottom: 4rpx solid transparent;
  white-space: nowrap;
}

.tab-item.active { color: #4a6cf7; border-bottom-color: #4a6cf7; font-weight: 600; }

.order-list { height: calc(100vh - 110rpx); padding: 24rpx; }

.order-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.04);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid #f3f4f6;
}

.order-shop { font-size: 28rpx; font-weight: 600; color: #1f2937; }

.order-status { font-size: 24rpx; font-weight: 600; }
.status--pending_payment { color: #f59e0b; }
.status--paid, .status--ready { color: #4a6cf7; }
.status--completed  { color: #10b981; }
.status--cancelled, .status--refunded { color: #9ca3af; }

.order-items { margin-bottom: 16rpx; }

.order-item { display: flex; align-items: center; padding: 10rpx 0; }

.item-img {
  width: 80rpx;
  height: 80rpx;
  border-radius: 10rpx;
  background: #f3f4f6;
  margin-right: 16rpx;
  flex-shrink: 0;
}

.item-detail { flex: 1; }
.item-name   { display: block; font-size: 26rpx; color: #1f2937; }
.item-qty    { display: block; font-size: 22rpx; color: #9ca3af; margin-top: 4rpx; }
.item-price  { font-size: 26rpx; color: #1f2937; font-weight: 500; }

.order-footer { padding-top: 16rpx; border-top: 1rpx solid #f3f4f6; text-align: right; }
.order-total  { font-size: 26rpx; color: #4b5563; }
.total-amount { font-size: 30rpx; font-weight: 700; color: #1f2937; }

/* 核销码 */
.verify-code-section {
  margin-top: 16rpx;
  padding: 20rpx;
  background: #f9fafb;
  border-radius: 12rpx;
  text-align: center;
}

.verify-label { display: block; font-size: 24rpx; color: #9ca3af; margin-bottom: 12rpx; }

.verify-code {
  font-size: 56rpx;
  font-weight: 700;
  color: #4a6cf7;
  letter-spacing: 8rpx;
}

/* 操作按钮 */
.order-actions { display: flex; gap: 16rpx; margin-top: 16rpx; justify-content: flex-end; }

.action-btn {
  padding: 14rpx 28rpx;
  border-radius: 100rpx;
  font-size: 24rpx;
  background: #f3f4f6;
  color: #4b5563;
  border: none;
}

.action-btn--primary { background: #4a6cf7; color: #fff; }
.action-btn--refund  { background: #fef2f2; color: #ef4444; }

/* 空状态 */
.empty-state { display: flex; flex-direction: column; align-items: center; padding-top: 100rpx; gap: 16rpx; }
.empty-icon  { font-size: 96rpx; }
.empty-text  { font-size: 28rpx; color: #9ca3af; }

.go-shop-btn {
  background: #4a6cf7;
  color: #fff;
  border-radius: 100rpx;
  font-size: 28rpx;
  padding: 16rpx 48rpx;
  border: none;
  margin-top: 12rpx;
}
</style>
