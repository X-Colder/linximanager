<template>
  <!-- 优惠券卡片 -->
  <view class="coupon-card" :class="{ 'coupon-card--used': status !== 'unused', 'coupon-card--selected': selected }">
    <view class="coupon-left">
      <text class="coupon-value">
        <text v-if="type === 'fixed'">¥<text class="coupon-amount">{{ value }}</text></text>
        <text v-else><text class="coupon-amount">{{ value * 10 }}</text>折</text>
      </text>
      <text class="coupon-condition">满{{ minAmount }}元可用</text>
    </view>

    <!-- 锯齿分割线 -->
    <view class="coupon-divider">
      <view class="divider-dot divider-dot--top" />
      <view class="divider-line" />
      <view class="divider-dot divider-dot--bottom" />
    </view>

    <view class="coupon-right">
      <text class="coupon-name">{{ name }}</text>
      <text class="coupon-expire">有效期至 {{ formatDate(endAt) }}</text>
      <view v-if="status === 'unused'" class="coupon-action">
        <button
          v-if="showUseBtn"
          class="use-btn"
          @tap.stop="emit('use')"
        >立即使用</button>
        <button
          v-else-if="showSelectBtn"
          class="select-btn"
          :class="{ selected }"
          @tap.stop="emit('select')"
        >{{ selected ? '已选择' : '使用' }}</button>
      </view>
      <text v-else class="coupon-status-text">{{ STATUS_MAP[status] }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
defineProps<{
  name: string
  type: 'fixed' | 'percent'
  value: number
  minAmount: number
  endAt: string
  status: 'unused' | 'used' | 'expired'
  showUseBtn?: boolean
  showSelectBtn?: boolean
  selected?: boolean
}>()

const emit = defineEmits<{
  use: []
  select: []
}>()

const STATUS_MAP = { unused: '未使用', used: '已使用', expired: '已过期' }

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return `${d.getFullYear()}.${String(d.getMonth() + 1).padStart(2, '0')}.${String(d.getDate()).padStart(2, '0')}`
}
</script>

<style scoped>
.coupon-card {
  display: flex;
  background: #fff;
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.06);
  border: 2rpx solid transparent;
  transition: border-color 0.2s;
}

.coupon-card--selected { border-color: #4a6cf7; }
.coupon-card--used     { opacity: 0.5; }

.coupon-left {
  background: linear-gradient(135deg, #4a6cf7, #3451d1);
  padding: 28rpx 24rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 180rpx;
}

.coupon-value { color: rgba(255,255,255,0.8); font-size: 24rpx; }
.coupon-amount { font-size: 52rpx; font-weight: 700; color: #fff; }
.coupon-condition { font-size: 20rpx; color: rgba(255,255,255,0.7); margin-top: 6rpx; }

.coupon-divider {
  position: relative;
  width: 20rpx;
  background: repeating-linear-gradient(
    to bottom,
    transparent 0,
    transparent 8rpx,
    #f3f4f6 8rpx,
    #f3f4f6 16rpx
  );
}

.divider-dot {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 28rpx;
  height: 28rpx;
  border-radius: 50%;
  background: #f9fafb;
}

.divider-dot--top    { top: -14rpx; }
.divider-dot--bottom { bottom: -14rpx; }

.coupon-right {
  flex: 1;
  padding: 24rpx 28rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.coupon-name   { font-size: 28rpx; font-weight: 600; color: #1f2937; }
.coupon-expire { font-size: 22rpx; color: #9ca3af; margin-top: 8rpx; }

.coupon-action { margin-top: 12rpx; }

.use-btn, .select-btn {
  display: inline-block;
  padding: 8rpx 24rpx;
  border-radius: 100rpx;
  font-size: 24rpx;
  border: none;
  background: #4a6cf7;
  color: #fff;
}

.select-btn.selected { background: #e8edff; color: #4a6cf7; }

.coupon-status-text {
  font-size: 24rpx;
  color: #9ca3af;
  margin-top: 12rpx;
}
</style>
