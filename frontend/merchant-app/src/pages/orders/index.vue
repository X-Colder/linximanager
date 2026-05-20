<template>
  <view class="page">
    <!-- 筛选栏 -->
    <view class="filter-bar">
      <view
        v-for="status in statusList"
        :key="status.value"
        class="filter-item"
        :class="{ active: currentStatus === status.value }"
        @tap="currentStatus = status.value"
      >
        {{ status.label }}
        <text v-if="status.count" class="filter-count">{{ status.count }}</text>
      </view>
    </view>

    <scroll-view scroll-y class="order-list" v-if="filteredOrders.length">
      <view v-for="order in filteredOrders" :key="order.id" class="order-card">
        <view class="order-header">
          <text class="order-no">订单 {{ order.order_no }}</text>
          <text class="order-status" :class="`status--${order.status}`">
            {{ STATUS_LABEL[order.status] }}
          </text>
        </view>
        <view class="order-items">
          <view v-for="item in order.items" :key="item.product_id" class="order-item">
            <text class="item-name">{{ item.product_name }}</text>
            <text class="item-qty">x{{ item.quantity }}</text>
            <text class="item-price">¥{{ item.total_price }}</text>
          </view>
        </view>
        <view class="order-footer">
          <text class="order-time">{{ formatTime(order.created_at) }}</text>
          <text class="order-total">共 ¥{{ order.pay_amount }}</text>
        </view>
        <view v-if="order.status === 'paid'" class="order-actions">
          <button class="action-btn action-btn--primary" @tap="verify(order)">扫码核销</button>
        </view>
      </view>
    </scroll-view>

    <view v-else class="empty-state">
      <text class="empty-text">暂无订单</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import http from '@/api/request'

type OrderStatus = 'pending_payment' | 'paid' | 'preparing' | 'ready' | 'completed' | 'cancelled'

interface OrderItem {
  product_id: number
  product_name: string
  quantity: number
  unit_price: number
  total_price: number
}

interface Order {
  id: number
  order_no: string
  status: OrderStatus
  pay_amount: number
  items: OrderItem[]
  created_at: string
}

const STATUS_LABEL: Record<OrderStatus, string> = {
  pending_payment: '待付款',
  paid: '待备货',
  preparing: '备货中',
  ready: '待取餐',
  completed: '已完成',
  cancelled: '已取消'
}

const statusList = [
  { label: '全部', value: '', count: 0 },
  { label: '待处理', value: 'paid', count: 0 },
  { label: '已完成', value: 'completed', count: 0 }
]

const currentStatus = ref('')
const orders = ref<Order[]>([])

const filteredOrders = computed(() => {
  if (!currentStatus.value) return orders.value
  return orders.value.filter(o => o.status === currentStatus.value)
})

function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN', { month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function verify(order: Order) {
  uni.navigateTo({ url: `/pages/verify/index?orderId=${order.id}` })
}

onMounted(async () => {
  try {
    orders.value = await http.get<Order[]>('/merchant/orders')
  } catch {
    orders.value = [
      {
        id: 1, order_no: 'LX202405201001', status: 'paid', pay_amount: 74,
        items: [
          { product_id: 1, product_name: '烤鸡翅(5串)', quantity: 2, unit_price: 15, total_price: 30 },
          { product_id: 2, product_name: '青岛啤酒', quantity: 4, unit_price: 6, total_price: 24 }
        ],
        created_at: new Date().toISOString()
      }
    ]
  }
})
</script>

<style scoped>
.page { min-height: 100vh; background: #f9fafb; }

.filter-bar {
  display: flex;
  background: #fff;
  padding: 0 24rpx;
  border-bottom: 1rpx solid #f3f4f6;
}

.filter-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 28rpx;
  color: #9ca3af;
  border-bottom: 4rpx solid transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
}

.filter-item.active { color: #4a6cf7; border-bottom-color: #4a6cf7; font-weight: 600; }

.filter-count {
  background: #ef4444;
  color: #fff;
  font-size: 20rpx;
  border-radius: 100rpx;
  padding: 2rpx 10rpx;
  min-width: 32rpx;
  text-align: center;
}

.order-list { height: calc(100vh - 120rpx); padding: 24rpx; }

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
  margin-bottom: 16rpx;
}

.order-no { font-size: 24rpx; color: #9ca3af; }

.order-status { font-size: 26rpx; font-weight: 600; }
.status--paid       { color: #f59e0b; }
.status--completed  { color: #10b981; }
.status--cancelled  { color: #9ca3af; }
.status--preparing  { color: #4a6cf7; }

.order-items { border-top: 1rpx solid #f3f4f6; padding-top: 16rpx; }

.order-item {
  display: flex;
  align-items: center;
  padding: 8rpx 0;
  font-size: 26rpx;
}

.item-name { flex: 1; color: #1f2937; }
.item-qty  { color: #9ca3af; margin: 0 24rpx; }
.item-price { color: #1f2937; font-weight: 500; }

.order-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f3f4f6;
}

.order-time  { font-size: 24rpx; color: #9ca3af; }
.order-total { font-size: 28rpx; font-weight: 700; color: #1f2937; }

.order-actions { margin-top: 16rpx; display: flex; gap: 16rpx; }

.action-btn {
  flex: 1;
  height: 72rpx;
  border-radius: 12rpx;
  font-size: 26rpx;
  border: none;
}

.action-btn--primary { background: #4a6cf7; color: #fff; }

.empty-state { display: flex; justify-content: center; padding: 80rpx 0; }
.empty-text  { font-size: 28rpx; color: #9ca3af; }
</style>
